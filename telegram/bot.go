package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"heytobi.dev/fuse/bot"

	"github.com/pkg/errors"
)

const (
	UpdateMethodWebhook    = "webhook"
	UpdateMethodGetUpdates = "getUpdates"

	telegramApiUrlFmt = "https://api.telegram.org/bot%s/%s"

	// Bot api methods
	//getUpdates  = "getUpdates"
	setWebhook  = "setWebhook"
	sendMessage = "sendMessage"

	httpPost = "POST"
)

var (
	errMissingToken        = errors.New("missing API token")
	errMissingWebhookUrl   = errors.New("a url is required to register a webhook")
	errMissingConfig       = errors.New("a configuration object is required to initialize a Bot connection")
	errMissingHttpClient   = errors.New("an http client is required to initialize a Bot connection")
	errHandlerExists       = errors.New("an handler already exists for this command")
	errInvalidUpdateMethod = errors.New("invalid update method")
	errInvalidUpdateType   = errors.New("invalid type passed to ProcessUpdate function")
	errInvalidSendableType = errors.New("invalid type passed to Send function")
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config defines the configurable parameters of a Bot.
type Config struct {
	Token          string
	UpdateMethod   string
	PollingTimeout int
}

// Bot defines the attributes of a Telegram Bot.
type Bot struct {
	config     *Config
	httpClient httpClient
	handlers   map[string]bot.HandlerFunc
}

// Init initializes a Bot instance.
//
// If no UpdateMethod is specified in the config, it defaults to getUpdates.
// See https://core.telegram.org/bots/api#getting-updates.
// Note that it only defaults to getUpdates if no update method is specified, if an invalid one is configured,
// an error is returned.
//
// It returns an error if any of these conditions are met:
//  - The given config is nil
//  - The configured UpdateMethod is invalid.
//  - The given serviceProvider is nil
func Init(config *Config, httpClient httpClient) (*Bot, error) {
	if config == nil {
		return nil, errMissingConfig
	}

	if config.Token == "" {
		return nil, errMissingToken
	}

	if config.UpdateMethod == "" {
		config.UpdateMethod = UpdateMethodGetUpdates
	}

	if config.UpdateMethod != UpdateMethodGetUpdates && config.UpdateMethod != UpdateMethodWebhook {
		return nil, errInvalidUpdateMethod
	}

	if httpClient == nil {
		return nil, errMissingHttpClient
	}

	bot := &Bot{
		config:     config,
		httpClient: httpClient,
		handlers:   make(map[string]bot.HandlerFunc),
	}

	return bot, nil
}

// Start starts the process of polling for updates from Bot.
func (b *Bot) Start() error {
	//if b.config.UpdateMethod == UpdateMethodGetUpdates {
	// start polling here
	//}

	return nil
}

// RegisterHandler registers the given handler function to handle invocations of the given command.
func (b *Bot) RegisterHandler(command string, handlerFunc bot.HandlerFunc) error {
	if _, handlerExists := b.handlers[command]; handlerExists {
		return errHandlerExists
	}

	b.handlers[command] = handlerFunc

	return nil
}

// RegisterWebhook registers the given webhook to listen for updates.
// Returns the result of the request, True on success.
// See https://core.telegram.org/bots/api#setwebhook
func (b *Bot) RegisterWebhook(url string) (bool, error) {
	if url == "" {
		return false, errMissingWebhookUrl
	}

	url = fmt.Sprintf(telegramApiUrlFmt, b.config.Token, setWebhook)
	requestBody := setWebhookRequest{
		Url: url,
	}

	responseBody, err := b.makeRequest(httpPost, url, requestBody)
	if err != nil {
		return false, errors.Wrap(err, "failed to make setWebhook request")
	}

	var resp setWebhookResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall setWebhook response")
	}

	return resp.Ok, nil
}

// ProcessUpdates processes updates from telegram.
func (b *Bot) ProcessUpdate(u bot.Update) error {
	update, isUpdate := u.(Update)
	if !isUpdate {
		return errInvalidUpdateType
	}

	if update.Message != nil {
		if handler, hasHandler := b.handlers[update.Message.Text]; hasHandler {
			handler(update)
		}
	}

	return nil
}

// Send sends a message to the user.
func (b *Bot) Send(s bot.Sendable) (bool, error) {
	message, isMessage := s.(SendMessageRequest)
	if !isMessage {
		return false, errInvalidSendableType
	}

	url := fmt.Sprintf(telegramApiUrlFmt, b.config.Token, sendMessage)
	responseBody, err := b.makeRequest(httpPost, url, message)
	if err != nil {
		return false, errors.Wrap(err, "failed to make sendMessage request")
	}

	var resp sendMessageResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall sendMessage response")
	}

	return resp.Ok, nil
}

func (b *Bot) makeRequest(method, url string, body interface{}) ([]byte, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	response, err := b.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse response body")
	}

	return responseBody, nil
}
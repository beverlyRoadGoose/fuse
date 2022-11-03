package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	UpdateMethodWebhook    = "webhook"
	UpdateMethodGetUpdates = "getUpdates"

	defaultBotApiServer = "https://api.telegram.org"

	httpPost = "POST"
)

var (
	errEmptyCommand            = errors.New("empty command")
	errNilUpdate               = errors.New("update cannot be nil")
	errNilMessageRequest       = errors.New("message cannot be nil")
	errMissingToken            = errors.New("missing API token")
	errMissingWebhookUrl       = errors.New("a url is required to register a webhook")
	errNilHttpClient           = errors.New("an http client is required to initialize a Bot connection")
	errNilPoller               = errors.New("a poller is required when using getUpdates")
	errNilConfig               = errors.New("a configuration object is required to initialize a Bot connection")
	errHandlerExists           = errors.New("an handler already exists for this command")
	errInvalidUpdateMethod     = errors.New("invalid update method")
	errDefaultHandlerExists    = errors.New("a default handler is already registered")
	errWrongUpdateMethodConfig = errors.New("bot is not configured to use webhook update method")
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type poller interface {
	start() error
	getUpdatesChannel() <-chan *Update
}

// Handler defines structs that can handle bot commands / messages.
type Handler interface {
	Handle(update *Update)
}

// HandlerFunc defines functions that can handle bot commands / messages.
type HandlerFunc func(update *Update)

// Config defines Telegrams configurable parameters.
type Config struct {
	BotApiServer        string
	BotApiServerPort    int
	Token               string
	UpdateMethod        string
	PollingCronSchedule string
	PollingCronTimezone string
	PollingTimeout      int
	PollingUpdatesLimit int
	AllowedUpdates      []string `json:"allowed_updates"`
}

// Bot defines the attributes of a Telegram Bot.
type Bot struct {
	config         *Config
	httpClient     httpClient
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
	poller         poller
	isRunning      bool
	apiUrlFmt      string
}

// NewBot initializes a Bot instance.
//
// If no UpdateMethod is specified in the config, it defaults to getUpdates.
// See https://core.telegram.org/bots/api#getting-updates.
// Note that it only defaults to getUpdates if no update method is specified, if an invalid one is configured,
// an error is returned.
//
// It returns an error if any of these conditions are met:
//   - The given config is nil
//   - The configured UpdateMethod is invalid.
func NewBot(config *Config, httpClient httpClient) (*Bot, error) {
	if config == nil {
		return nil, errNilConfig
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
		return nil, errNilHttpClient
	}

	bot := &Bot{
		config:     config,
		httpClient: httpClient,
		handlers:   make(map[string]HandlerFunc),
		apiUrlFmt:  deriveBotApiUrlBase(config) + "/bot%s/%s",
	}

	return bot, nil
}

// WithPoller sets the poller used to continuously query for updates
func (b *Bot) WithPoller(p poller) *Bot {
	b.poller = p
	return b
}

// Start starts the process of polling for updates from Telegram.
func (b *Bot) Start() error {
	if b.config.UpdateMethod == UpdateMethodGetUpdates && !b.isRunning {
		if b.poller == nil {
			return errNilPoller
		}

		updatesChan := b.poller.getUpdatesChannel()
		go func() {
			for update := range updatesChan {
				err := b.ProcessUpdate(update)
				if err != nil {
					logrus.WithError(err).Error("failed to process update")
				}
			}
		}()

		_, err := b.deleteWebhook(false)
		if err != nil {
			return errors.Wrap(err, "failed to delete webhook")
		}

		err = b.poller.start()
		if err != nil {
			return errors.Wrap(err, "failed to start poller")
		}

		b.isRunning = true
	}

	return nil
}

// RegisterWebhook registers the given webhook to listen for updates.
// Returns the result of the request, True on success.
// See https://core.telegram.org/bots/api#setwebhook
func (b *Bot) RegisterWebhook(webhook *Webhook) (bool, error) {
	if b.config.UpdateMethod == UpdateMethodGetUpdates {
		return false, errWrongUpdateMethodConfig
	}

	if webhook.Url == "" {
		return false, errMissingWebhookUrl
	}

	url := fmt.Sprintf(b.apiUrlFmt, b.config.Token, endpointSetWebhook)

	if webhook.AllowedUpdates == nil {
		webhook.AllowedUpdates = b.config.AllowedUpdates
	}

	bodyJson, err := json.Marshal(webhook)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal register webhook request body")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return false, errors.Wrap(err, "failed to create register webhook request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := b.httpClient.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "register webhook http request failed")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse register webhook response body")
	}

	var resp webhookResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall setWebhook response")
	}

	return resp.Ok, nil
}

// deleteWebhook deletes the registered webhook.
// See https://core.telegram.org/bots/api#deletewebhook
func (b *Bot) deleteWebhook(dropPendingUpdates bool) (bool, error) {
	url := fmt.Sprintf(b.apiUrlFmt, b.config.Token, endpointDeleteWebhook)

	body := deleteWebhookRequest{DropPendingUpdates: dropPendingUpdates}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal delete webhook request body")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return false, errors.Wrap(err, "failed to create register webhook request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := b.httpClient.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "delete webhook http request failed")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse delete webhook response body")
	}

	var resp webhookResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall delete Webhook response")
	}

	return resp.Ok, nil
}

// RegisterDefaultHandler registers the given handler function as the default. The default handler handles all
// updates that don't match a specific command that is assigned its own handler in RegisterHandler.
func (b *Bot) RegisterDefaultHandler(handler HandlerFunc) error {
	if b.defaultHandler != nil {
		return errDefaultHandlerExists
	}

	b.defaultHandler = handler

	return nil
}

// RegisterHandler registers the given handler function to handle invocations of the given command.
func (b *Bot) RegisterHandler(command string, handlerFunc HandlerFunc) error {
	if command == "" {
		return errEmptyCommand
	}

	if _, handlerExists := b.handlers[command]; handlerExists {
		return errHandlerExists
	}

	b.handlers[command] = handlerFunc

	return nil
}

// ProcessUpdates processes updates from telegram.
func (b *Bot) ProcessUpdate(update *Update) error {
	if update == nil {
		return errNilUpdate
	}

	if update.Message != nil {
		if handler, hasHandler := b.handlers[update.Message.Text]; hasHandler {
			if handler != nil {
				handler(update)
			}
			return nil
		}

		if b.defaultHandler != nil {
			b.defaultHandler(update)
		}
	}

	return nil
}

// SendMessage sends a message to the user.
func (b *Bot) SendMessage(message *SendMessageRequest) (bool, error) {
	if message == nil {
		return false, errNilMessageRequest
	}

	url := fmt.Sprintf(b.apiUrlFmt, b.config.Token, endpointSendMessage)

	bodyJson, err := json.Marshal(message)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal send request")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return false, errors.Wrap(err, "failed to create send request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := b.httpClient.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "http request failed")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse send response body")
	}

	var resp sendMessageResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall sendMessage response")
	}

	return resp.Ok, nil
}

func deriveBotApiUrlBase(config *Config) string {
	botApiUrlBase := defaultBotApiServer
	if config.BotApiServer != "" {
		botApiUrlBase = config.BotApiServer
		if config.BotApiServerPort != 0 {
			botApiUrlBase += ":" + strconv.Itoa(config.BotApiServerPort)
		}
	}

	return botApiUrlBase
}

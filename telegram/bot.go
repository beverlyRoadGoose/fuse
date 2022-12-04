package telegram // import "heytobi.dev/fuse/telegram"

import (
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

type ActionResult struct {
	Successful  bool
	Description string
}

// Handler defines structs that can handle bot commands / messages.
type Handler interface {
	Handle(update *Update) error
}

// HandlerFunc defines functions that can handle bot commands / messages.
type HandlerFunc func(update *Update) error

// Config defines Telegrams configurable parameters.
type Config struct {
	BotApiServer        string
	BotApiServerPort    int
	Token               string
	UpdateMethod        string
	PollingIntervalMS   int64
	PollingTimeout      int
	PollingUpdatesLimit int
	AllowedUpdates      []string `json:"allowed_updates"`
}

// Bot defines the attributes of a Telegram Bot.
type Bot struct {
	config           *Config
	httpClient       httpClient
	handlers         map[string]HandlerFunc
	defaultHandler   HandlerFunc
	poller           poller
	isRunning        bool
	apiUrlFmt        string
	messagingService *messagingService
	webhookService   *webhookService
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

	apiUrlFmt := deriveBotApiUrlBase(config) + "/bot%s/%s"

	webhookService, err := newWebhookService(httpClient, apiUrlFmt, config.Token, config.AllowedUpdates)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize webhook service")
	}

	messagingService, err := newMessagingService(httpClient, apiUrlFmt, config.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize messaging service")
	}

	bot := &Bot{
		config:           config,
		httpClient:       httpClient,
		handlers:         make(map[string]HandlerFunc),
		apiUrlFmt:        apiUrlFmt,
		messagingService: messagingService,
		webhookService:   webhookService,
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

		_, err := b.webhookService.deleteWebhook(false)
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
	return b.webhookService.registerWebhook(webhook)
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

// ProcessUpdate processes updates from telegram.
func (b *Bot) ProcessUpdate(update *Update) error {
	if update == nil {
		return errNilUpdate
	}

	if update.Message != nil {
		if handler, hasHandler := b.handlers[update.Message.Text]; hasHandler {
			if handler != nil {
				err := handler(update)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if b.defaultHandler != nil {
			err := b.defaultHandler(update)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SendMessage sends a message to the user.
func (b *Bot) SendMessage(message *SendMessageRequest) (*ActionResult, error) {
	return b.messagingService.sendMessage(message)
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

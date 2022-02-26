package telegram // import "heytobi.dev/fuse/telegram"

import (
	"net/http"

	"heytobi.dev/fuse/bot"

	"github.com/pkg/errors"
)

const (
	UpdateMethodWebhook    = "webhook"
	UpdateMethodGetUpdates = "getUpdates"
)

var (
	errMissingToken        = errors.New("missing API token")
	errMissingConfig       = errors.New("a configuration object is required to initialize a Telegram connection")
	errMissingHttpClient   = errors.New("an http client is required to initialize a Telegram connection")
	errInvalidUpdateMethod = errors.New("invalid update method")
	errHandlerExists       = errors.New("an handler already exists for this command")
	errInvalidUpdateType   = errors.New("invalid type passed to ProcessUpdate function")
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config defines the configurable parameters of a Telegram connection.
type Config struct {
	Token          string
	UpdateMethod   string
	PollingTimeout int
}

// Telegram defines an instance of telegram connection.
type Telegram struct {
	config     *Config
	httpClient httpClient
	handlers   map[string]bot.HandlerFunc
}

// Init initializes a Telegram instance.
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
func Init(config *Config, httpClient httpClient) (*Telegram, error) {
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

	telegram := &Telegram{
		config:     config,
		httpClient: httpClient,
		handlers:   make(map[string]bot.HandlerFunc),
	}

	return telegram, nil
}

// Start starts the process of polling for updates from Telegram.
func (t *Telegram) Start() error {
	return nil
}

// Send sends a message to the user.
func (t *Telegram) Send(message bot.Sendable) error {
	return nil
}

// RegisterHandler registers the given handler function to handle invocations of the given command.
func (t *Telegram) RegisterHandler(command string, handlerFunc bot.HandlerFunc) error {
	if _, handlerExists := t.handlers[command]; handlerExists {
		return errHandlerExists
	}

	t.handlers[command] = handlerFunc

	return nil
}

// ProcessUpdates processes updates from telegram.
func (t *Telegram) ProcessUpdate(update bot.Update) error {
	update, ok := update.(Update)
	if !ok {
		return errInvalidUpdateType
	}

	return nil
}

// registerWebhook registers the configured webhook to listen for events.
func (t *Telegram) registerWebhook() error {
	return nil
}

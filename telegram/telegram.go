package telegram // import "heytobi.dev/fuse/telegram"

import (
	"net/http"

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
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config defines the configurable parameters of a Telegram connection.
type Config struct {
	Token        string
	UpdateMethod string
}

// Telegram defines an instance of telegram connection.
type Telegram struct {
	config     *Config
	httpClient httpClient
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
	}

	return telegram, nil
}

// SendMessage sends a message to the user.
func (t *Telegram) SendMessage() error {
	return nil
}

// Start starts the process of polling for updates from Telegram.
func (t *Telegram) Start() error {
	return nil
}

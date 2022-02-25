package telegram // import "heytobi.dev/fuse/telegram"

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	errMissingToken      = errors.New("missing API token")
	errMissingConfig     = errors.New("a configuration object is required to initialize a Telegram connection")
	errMissingHttpClient = errors.New("an http client is required to initialize a Telegram connection")
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config defines the configurable parameters of a Telegram connection.
type Config struct {
	token string
}

// Telegram defines an instance of telegram connection.
type Telegram struct {
	config     *Config
	httpClient httpClient
}

// Init initializes a Telegram instance.
//
// It returns an error if any of these conditions are met:
//  - The given config is nil
//  - The given serviceProvider is nil
func Init(config *Config, httpClient httpClient) (*Telegram, error) {
	if config == nil {
		return nil, errMissingConfig
	}

	if config.token == "" {
		return nil, errMissingToken
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

// SendMessage sends a message to the user
func (t *Telegram) SendMessage() error {
	return nil
}

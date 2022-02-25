package bot // import "heytobi.dev/fuse/bot"

import (
	"github.com/pkg/errors"
)

var (
	errMissingConfig          = errors.New("a configuration object is required to initialize a bot")
	errMissingServiceProvider = errors.New("a service provider is required to initialize a bot")
)

// messagingServiceProvider defines the functions required of a messaging service provider.
type messagingServiceProvider interface {
	SendMessage() error
}

// Config defines the configurable parameters of a Bot.
type Config struct {
}

// Bot defines the attributes of a bot.
type Bot struct {
	config          *Config
	serviceProvider messagingServiceProvider
}

// NewBot initializes a new bot.
//
// It returns an error if any of these conditions are met:
// - The given config is nil
// - The given serviceProvider is nil
func NewBot(config *Config, serviceProvider messagingServiceProvider) (*Bot, error) {
	if config == nil {
		return nil, errMissingConfig
	}

	if serviceProvider == nil {
		return nil, errMissingServiceProvider
	}

	bot := &Bot{
		config:          config,
		serviceProvider: serviceProvider,
	}

	return bot, nil
}

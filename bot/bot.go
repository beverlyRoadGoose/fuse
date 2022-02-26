package bot // import "heytobi.dev/fuse/bot"

import (
	"github.com/pkg/errors"
)

var (
	errMissingServiceProvider = errors.New("a service provider is required to initialize a bot")
)

// HandlerFunc defines functions that can handle bot commands / messages.
type HandlerFunc func(interface{})

// Sendable defines sendable items.
type Sendable interface{}

// messagingServiceProvider defines the functions required of a messaging service provider.
type messagingServiceProvider interface {
	Start() error
	Send(message Sendable) error
	RegisterHandler(command string, handlerFunc HandlerFunc) error
}

// Bot defines the attributes of a bot.
type Bot struct {
	serviceProvider messagingServiceProvider
}

// NewBot initializes a new bot.
//
// It returns an error if any of these conditions are met:
//  - The given serviceProvider is nil
func NewBot(serviceProvider messagingServiceProvider) (*Bot, error) {
	if serviceProvider == nil {
		return nil, errMissingServiceProvider
	}

	bot := &Bot{
		serviceProvider: serviceProvider,
	}

	return bot, nil
}

// Start starts the process of polling for updates from the service provider.
func (b *Bot) Start() error {
	return b.serviceProvider.Start()
}

// Send sends a message to the user.
func (b *Bot) Send(message Sendable) error {
	return b.serviceProvider.Send(message)
}

// RegisterHandler registers the given handler function to handle invocations of the given command.
func (b *Bot) RegisterHandler(command string, handlerFunc HandlerFunc) error {
	return b.serviceProvider.RegisterHandler(command, handlerFunc)
}

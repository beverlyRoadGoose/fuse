package bot // import "heytobi.dev/fuse/bot"

// HandlerFunc defines functions that can handle bot commands / messages.
type HandlerFunc func(interface{})

// Sendable defines sendable items.
type Sendable interface{}

// Update defines updates that can be received from messaging service providers.
type Update interface{}

// Bot defines the functions required of a messaging service provider.
type Bot interface {
	Start() error
	Send(message Sendable) (bool, error)
	RegisterHandler(command string, handler HandlerFunc) error
	RegisterWebhook(url string) (bool, error)
	ProcessUpdate(update Update) error
}

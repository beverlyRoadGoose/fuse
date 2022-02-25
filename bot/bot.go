package bot // import "heytobi.dev/fuse/bot"

// MessagingServiceProvider defines the functions required of a messaging service provider.
type MessagingServiceProvider interface {
	SendMessage() error
}

// Bot defines the attributes of a bot.
type Bot struct {
	serviceProvider MessagingServiceProvider
}

// New initializes a new bot.
func New(serviceProvider MessagingServiceProvider) (*Bot, error) {
	bot := &Bot{serviceProvider: serviceProvider}

	err := bot.serviceProvider.SendMessage()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

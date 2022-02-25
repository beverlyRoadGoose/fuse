package bot // import "heytobi.dev/fuse/bot"

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewBot_ReturnErrorIfConfigIsNil(t *testing.T) {
	bot, err := NewBot(nil, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingConfig)
}

func TestNewBot_ReturnErrorIfServiceProviderIsMissing(t *testing.T) {
	bot, err := NewBot(&Config{}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingServiceProvider)
}

func TestNewBot_InitializeBotSuccessfully(t *testing.T) {
	bot, err := NewBot(&Config{}, &mockMessagingServiceProvider{})

	assert.NotNil(t, bot)
	assert.Nil(t, err)
}

func TestBot_Start_ReturnErrorIfServiceProviderFailsToStart(t *testing.T) {
	serviceProvider := &mockMessagingServiceProvider{}
	serviceProvider.On("Start").Return(errors.New("failed to start"))

	bot, _ := NewBot(&Config{}, serviceProvider)
	err := bot.Start()

	assert.NotNil(t, err)
}

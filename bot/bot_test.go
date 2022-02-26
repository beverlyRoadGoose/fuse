package bot // import "heytobi.dev/fuse/bot"

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewBot_ReturnErrorIfServiceProviderIsMissing(t *testing.T) {
	bot, err := NewBot(nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errMissingServiceProvider, err)
}

func TestNewBot_InitializeBotSuccessfully(t *testing.T) {
	bot, err := NewBot(&mockMessagingServiceProvider{})

	assert.NotNil(t, bot)
	assert.Nil(t, err)
}

func TestBot_Start_ReturnErrorIfServiceProviderFailsToStart(t *testing.T) {
	serviceProvider := &mockMessagingServiceProvider{}
	serviceProvider.On("Start").Return(errors.New("failed to start"))

	bot, _ := NewBot(serviceProvider)
	err := bot.Start()

	assert.NotNil(t, err)
}

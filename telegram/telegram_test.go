package telegram // import "heytobi.dev/fuse/telegram"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_ReturnErrorIfConfigIsNil(t *testing.T) {
	bot, err := Init(nil, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errMissingConfig, err)
}

func TestInit_ReturnErrorIfTokenIsMissing(t *testing.T) {
	bot, err := Init(&Config{}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errMissingToken, err)
}

func TestInit_DefaultToGetUpdatesIfNoUpdateMethodIsSpecified(t *testing.T) {
	bot, err := Init(&Config{Token: "test"}, &mockHttpClient{})

	assert.Nil(t, err)
	assert.Equal(t, UpdateMethodGetUpdates, bot.config.UpdateMethod)
}

func TestInit_ReturnErrorIfUpdateMethodIsInvalid(t *testing.T) {
	bot, err := Init(
		&Config{
			Token:        "test",
			UpdateMethod: "invalid",
		},
		nil,
	)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errInvalidUpdateMethod, err)
}

func TestInit_ReturnErrorIfHttpClientIsMissing(t *testing.T) {
	bot, err := Init(&Config{Token: "test"}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errMissingHttpClient, err)
}

func TestInit_InitializeSuccessfully(t *testing.T) {
	telegram, err := Init(&Config{Token: "test"}, &mockHttpClient{})

	assert.NotNil(t, telegram)
	assert.Nil(t, err)
}

func TestStart_StartSuccessfully(t *testing.T) {
	telegram, err := Init(&Config{Token: "test"}, &mockHttpClient{})
	err = telegram.Start()

	assert.Nil(t, err)
}

func TestSend_SendMessageSuccessfully(t *testing.T) {
	telegram, err := Init(&Config{Token: "test"}, &mockHttpClient{})
	err = telegram.Send("message")

	assert.Nil(t, err)
}

func TestRegisterHandler_RegisterHandlerSuccessfully(t *testing.T) {
	telegram, err := Init(&Config{Token: "test"}, &mockHttpClient{})
	err = telegram.RegisterHandler("/start", func(message interface{}) {})

	assert.Nil(t, err)
}

func TestRegisterHandler_ReturnErrorIfHandlerExists(t *testing.T) {
	telegram, err := Init(&Config{Token: "test"}, &mockHttpClient{})
	err = telegram.RegisterHandler("/start", func(message interface{}) {})

	// try registering another handler for the same command
	err = telegram.RegisterHandler("/start", func(message interface{}) {})

	assert.NotNil(t, err)
	assert.Equal(t, errHandlerExists, err)
}

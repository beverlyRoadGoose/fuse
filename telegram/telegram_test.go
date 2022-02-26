package telegram // import "heytobi.dev/fuse/telegram"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_ReturnErrorIfConfigIsNil(t *testing.T) {
	bot, err := Init(nil, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingConfig)
}

func TestInit_ReturnErrorIfTokenIsMissing(t *testing.T) {
	bot, err := Init(&Config{}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingToken)
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
	assert.Equal(t, err, errInvalidUpdateMethod)
}

func TestInit_ReturnErrorIfHttpClientIsMissing(t *testing.T) {
	bot, err := Init(&Config{Token: "test"}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingHttpClient)
}

func TestInit_InitializeSuccessfully(t *testing.T) {
	telegram, err := Init(&Config{Token: "test"}, &mockHttpClient{})

	assert.NotNil(t, telegram)
	assert.Nil(t, err)
}

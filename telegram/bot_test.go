package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	bot, err := Init(&Config{Token: "test"}, &mockHttpClient{})

	assert.NotNil(t, bot)
	assert.Nil(t, err)
}

func TestInit_InitializePollerIfUpdateMethodIsNotWebhook(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})

	assert.NotNil(t, bot.poller)
}

func TestInit_DoNotInitializePollerIfUpdateMethodIsWebhook(t *testing.T) {
	bot, _ := Init(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, &mockHttpClient{})

	assert.Nil(t, bot.poller)
}

func TestStart_StartSuccessfully(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.Start()

	assert.Nil(t, err)
}

func TestSend_ReturnErrorIfMessageIsNil(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	result, err := bot.Send(nil)

	assert.Error(t, err)
	assert.Equal(t, errNilMessageRequest, err)
	assert.False(t, result)
}

func TestSend_SendSuccessfully(t *testing.T) {
	response := sendMessageResponse{Ok: true}
	responseJson, _ := json.Marshal(response)
	body := ioutil.NopCloser(bytes.NewBuffer(responseJson))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	bot, _ := Init(&Config{Token: "test"}, httpClient)
	result, err := bot.Send(&SendMessageRequest{
		ChatID: 0,
		Text:   "test",
	})

	assert.True(t, result)
	assert.Nil(t, err)
}

func TestRegisterHandler_RegisterHandlerSuccessfully(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.RegisterHandler("/start", func(update *Update) {})

	assert.Nil(t, err)
	assert.NotNil(t, bot.handlers["/start"])
}

func TestRegisterHandler_ReturnErrorIfHandlerExists(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	_ = bot.RegisterHandler("/start", func(update *Update) {})

	// try registering another handler for the same command
	err := bot.RegisterHandler("/start", func(update *Update) {})

	assert.NotNil(t, err)
	assert.Equal(t, errHandlerExists, err)
}

func TestRegisterDefaultHandler_RegisterHandlerSuccessfully(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.RegisterDefaultHandler(func(update *Update) {})

	assert.Nil(t, err)
	assert.NotNil(t, bot.defaultHandler)
}

func TestRegisterDefaultHandler_ReturnErrorIfHandlerExists(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	_ = bot.RegisterDefaultHandler(func(update *Update) {})
	err := bot.RegisterDefaultHandler(func(update *Update) {})

	assert.NotNil(t, err)
	assert.Equal(t, errDefaultHandlerExists, err)
}

func TestRegisterWebhook_ReturnErrorIfUrlIsEmpty(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("failed"))

	bot, _ := Init(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	registered, err := bot.RegisterWebhook(&Webhook{Url: ""})

	assert.NotNil(t, err)
	assert.Equal(t, errMissingWebhookUrl, err)
	assert.False(t, registered)
}

func TestRegisterWebhook_ReturnErrorIfUpdateMethodIsNotWebhook(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	registered, err := bot.RegisterWebhook(&Webhook{Url: "url.test"})

	assert.NotNil(t, err)
	assert.Equal(t, errWrongUpdateMethodConfig, err)
	assert.False(t, registered)
}

func TestRegisterWebhook_ReturnErrorIfApiRequestFails(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("failed"))

	bot, _ := Init(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	result, err := bot.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestRegisterWebhook_RegisterSuccessfully(t *testing.T) {
	response := setWebhookResponse{Ok: true}
	json, _ := json.Marshal(response)
	body := ioutil.NopCloser(bytes.NewBuffer(json))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	bot, _ := Init(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	result, err := bot.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.True(t, result)
	assert.Nil(t, err)
}

func TestRegisterWebhook_ReturnFalseIfResponseResultIsFalse(t *testing.T) {
	response := setWebhookResponse{Ok: false}
	responseJson, _ := json.Marshal(response)
	body := ioutil.NopCloser(bytes.NewBuffer(responseJson))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	bot, _ := Init(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	result, err := bot.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.False(t, result)
	assert.Nil(t, err)
}

func TestProcessUpdate_ReturnErrorIfUpdateIsNil(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.ProcessUpdate(nil)

	assert.Error(t, err)
	assert.Equal(t, errNilUpdate, err)
}

func TestProcessUpdate_DontReturnErrorIfGivenValidUpdateType(t *testing.T) {
	bot, _ := Init(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.ProcessUpdate(&Update{
		Message: &Message{Text: "/test"},
	})

	assert.Nil(t, err)
}

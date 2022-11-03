package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBot_ReturnErrorIfConfigIsNil(t *testing.T) {
	bot, err := NewBot(nil, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errNilConfig, err)
}

func TestNewBot_ReturnErrorIfTokenIsMissing(t *testing.T) {
	bot, err := NewBot(&Config{}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errMissingToken, err)
}

func TestNewBot_DefaultToGetUpdatesIfNoUpdateMethodIsSpecified(t *testing.T) {
	bot, err := NewBot(&Config{Token: "test"}, &mockHttpClient{})

	assert.NoError(t, err)
	assert.Equal(t, UpdateMethodGetUpdates, bot.config.UpdateMethod)
}

func TestNewBot_ReturnErrorIfUpdateMethodIsInvalid(t *testing.T) {
	bot, err := NewBot(
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

func TestNewBot_ReturnErrorIfHttpClientIsMissing(t *testing.T) {
	bot, err := NewBot(&Config{Token: "test"}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, errNilHttpClient, err)
}

func TestNewBot_InitializeSuccessfully(t *testing.T) {
	bot, err := NewBot(&Config{Token: "test"}, &mockHttpClient{})

	assert.NotNil(t, bot)
	assert.NoError(t, err)
}

func TestNewBot_InitializeWithDefaultBotApiServer(t *testing.T) {
	bot, err := NewBot(&Config{Token: "test"}, &mockHttpClient{})

	assert.NotNil(t, bot)
	assert.NoError(t, err)
	assert.Equal(t, "https://api.telegram.org/bot%s/%s", bot.apiUrlFmt)
}

func TestNewBot_InitializeWithCustomBotApiServer(t *testing.T) {
	customServer := "http://custom.server"
	bot, err := NewBot(&Config{Token: "test", BotApiServer: customServer}, &mockHttpClient{})

	assert.NotNil(t, bot)
	assert.NoError(t, err)
	assert.Equal(t, customServer+"/bot%s/%s", bot.apiUrlFmt)
}

func TestNewBot_InitializeWithCustomBotApiServerAndPort(t *testing.T) {
	customServer := "http://custom.server"
	customPort := 9090
	bot, err := NewBot(&Config{Token: "test", BotApiServer: customServer, BotApiServerPort: customPort}, &mockHttpClient{})

	assert.NotNil(t, bot)
	assert.NoError(t, err)

	expectedUrlFmt := customServer + ":" + strconv.Itoa(customPort) + "/bot%s/%s"
	assert.Equal(t, expectedUrlFmt, bot.apiUrlFmt)
}

func TestStart_ReturnErrorIfUsingGetUpdatesAndPollerIsNil(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test", UpdateMethod: UpdateMethodGetUpdates}, &mockHttpClient{})
	err := bot.Start()

	assert.Error(t, err)
	assert.Equal(t, errNilPoller, err)
}

func TestStart_StartSuccessfully(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, &mockHttpClient{})
	err := bot.Start()

	assert.NoError(t, err)
}

func TestSend_ReturnErrorIfMessageIsNil(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	result, err := bot.SendMessage(nil)

	assert.Error(t, err)
	assert.Equal(t, errNilMessageRequest, err)
	assert.False(t, result)
}

func TestSend_SendSuccessfully(t *testing.T) {
	response := sendMessageResponse{Ok: true}
	responseJson, _ := json.Marshal(response)
	body := io.NopCloser(bytes.NewBuffer(responseJson))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	bot, _ := NewBot(&Config{Token: "test"}, httpClient)
	result, err := bot.SendMessage(&SendMessageRequest{
		ChatID: 0,
		Text:   "test",
	})

	assert.True(t, result)
	assert.NoError(t, err)
}

func TestRegisterHandler_RegisterHandlerSuccessfully(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.RegisterHandler("/start", func(update *Update) {})

	assert.NoError(t, err)
	assert.NotNil(t, bot.handlers["/start"])
}

func TestRegisterHandler_ReturnErrorIfHandlerExists(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	_ = bot.RegisterHandler("/start", func(update *Update) {})

	// try registering another handler for the same command
	err := bot.RegisterHandler("/start", func(update *Update) {})

	assert.Error(t, err)
	assert.Equal(t, errHandlerExists, err)
}

func TestRegisterDefaultHandler_RegisterHandlerSuccessfully(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.RegisterDefaultHandler(func(update *Update) {})

	assert.NoError(t, err)
	assert.NotNil(t, bot.defaultHandler)
}

func TestRegisterDefaultHandler_ReturnErrorIfHandlerExists(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	_ = bot.RegisterDefaultHandler(func(update *Update) {})
	err := bot.RegisterDefaultHandler(func(update *Update) {})

	assert.Error(t, err)
	assert.Equal(t, errDefaultHandlerExists, err)
}

func TestRegisterWebhook_ReturnErrorIfUrlIsEmpty(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("failed"))

	bot, _ := NewBot(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	registered, err := bot.RegisterWebhook(&Webhook{Url: ""})

	assert.Error(t, err)
	assert.Equal(t, errMissingWebhookUrl, err)
	assert.False(t, registered)
}

func TestRegisterWebhook_ReturnErrorIfUpdateMethodIsNotWebhook(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	registered, err := bot.RegisterWebhook(&Webhook{Url: "url.test"})

	assert.Error(t, err)
	assert.Equal(t, errWrongUpdateMethodConfig, err)
	assert.False(t, registered)
}

func TestRegisterWebhook_ReturnErrorIfApiRequestFails(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("failed"))

	bot, _ := NewBot(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	result, err := bot.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.Error(t, err)
	assert.False(t, result)
}

func TestRegisterWebhook_RegisterSuccessfully(t *testing.T) {
	response := webhookResponse{Ok: true}
	json, _ := json.Marshal(response)
	body := io.NopCloser(bytes.NewBuffer(json))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	bot, _ := NewBot(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	result, err := bot.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.True(t, result)
	assert.NoError(t, err)
}

func TestRegisterWebhook_ReturnFalseIfResponseResultIsFalse(t *testing.T) {
	response := webhookResponse{Ok: false}
	responseJson, _ := json.Marshal(response)
	body := io.NopCloser(bytes.NewBuffer(responseJson))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	bot, _ := NewBot(&Config{Token: "test", UpdateMethod: UpdateMethodWebhook}, httpClient)
	result, err := bot.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.False(t, result)
	assert.NoError(t, err)
}

func TestProcessUpdate_ReturnErrorIfUpdateIsNil(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.ProcessUpdate(nil)

	assert.Error(t, err)
	assert.Equal(t, errNilUpdate, err)
}

func TestProcessUpdate_DontReturnErrorIfGivenValidUpdateType(t *testing.T) {
	bot, _ := NewBot(&Config{Token: "test"}, &mockHttpClient{})
	err := bot.ProcessUpdate(&Update{
		Message: &Message{Text: "/test"},
	})

	assert.NoError(t, err)
}

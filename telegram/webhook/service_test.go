package webhook // import "heytobi.dev/fuse/telegram/webhook

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const apiUrlFmt = "test.com/bot%s/%s"

func TestRegisterWebhook_ReturnErrorIfUrlIsEmpty(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("failed"))

	service, _ := NewService(httpClient, "", "test", nil)
	registered, err := service.RegisterWebhook(&Webhook{Url: ""})

	assert.Error(t, err)
	assert.Equal(t, errMissingWebhookUrl, err)
	assert.False(t, registered)
}

func TestRegisterWebhook_ReturnErrorIfApiRequestFails(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("failed"))

	service, _ := NewService(httpClient, apiUrlFmt, "test", nil)
	result, err := service.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.Error(t, err)
	assert.False(t, result)
}

func TestRegisterWebhook_RegisterSuccessfully(t *testing.T) {
	response := webhookResponse{Ok: true}
	json, _ := json.Marshal(response)
	body := io.NopCloser(bytes.NewBuffer(json))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	service, _ := NewService(httpClient, apiUrlFmt, "test", nil)
	result, err := service.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.True(t, result)
	assert.NoError(t, err)
}

func TestRegisterWebhook_ReturnFalseIfResponseResultIsFalse(t *testing.T) {
	response := webhookResponse{Ok: false}
	responseJson, _ := json.Marshal(response)
	body := io.NopCloser(bytes.NewBuffer(responseJson))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	service, _ := NewService(httpClient, apiUrlFmt, "test", nil)
	result, err := service.RegisterWebhook(&Webhook{Url: "webhook.url"})

	assert.False(t, result)
	assert.NoError(t, err)
}

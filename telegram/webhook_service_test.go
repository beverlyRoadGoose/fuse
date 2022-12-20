package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteWebhook_CanDeleteWebhook(t *testing.T) {
	httpClient := &mockHttpClient{}
	response := webhookResponse{Ok: true}
	responseJson, _ := json.Marshal(response)
	responseBody := io.NopCloser(bytes.NewBufferString(string(responseJson)))
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       responseBody,
	}, nil)

	var allowedUpdates []string
	service, _ := newWebhookService(httpClient, testApiUrlFmt, testToken, allowedUpdates)

	success, err := service.deleteWebhook(true)

	assert.True(t, success)
	assert.Nil(t, err)
}

func TestDeleteWebhook_ReturnsFalseIfResponseCodeIsNot200(t *testing.T) {
	httpClient := &mockHttpClient{}
	response := webhookResponse{Ok: false}
	responseJson, _ := json.Marshal(response)
	responseBody := io.NopCloser(bytes.NewBufferString(string(responseJson)))
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       responseBody,
	}, nil)

	var allowedUpdates []string
	service, _ := newWebhookService(httpClient, testApiUrlFmt, testToken, allowedUpdates)

	success, err := service.deleteWebhook(true)

	assert.False(t, success)
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to delete webhook, unexpected response code"))
}

func TestDeleteWebhook_ReturnsFalseIfDeleteWebhookRequestFails(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	var allowedUpdates []string
	service, _ := newWebhookService(httpClient, testApiUrlFmt, testToken, allowedUpdates)

	success, err := service.deleteWebhook(true)

	assert.False(t, success)
	assert.Error(t, err)
	assert.True(t, strings.EqualFold(err.Error(), "delete webhook http request failed: error"))
}

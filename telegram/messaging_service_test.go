package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testApiUrlFmt = "https://api.telegram.local/bot%s/%s"
	testToken     = ""
)

func TestMessagingService_ReturnsErrorIfMessageIsNil(t *testing.T) {
	httpClient := &mockHttpClient{}
	service, _ := newMessagingService(httpClient, testApiUrlFmt, testToken)

	result, err := service.sendMessage(nil)

	assert.False(t, result.Successful)
	assert.Error(t, err)
	assert.Equal(t, errNilMessageRequest, err)
}

func TestMessagingService_ActionResultIsSuccessfulIfMessageIsSentSuccessfully(t *testing.T) {
	httpClient := &mockHttpClient{}
	response := sendMessageResponse{Ok: true}
	responseJson, _ := json.Marshal(response)
	responseBody := io.NopCloser(bytes.NewBufferString(string(responseJson)))
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       responseBody,
	}, nil)
	service, _ := newMessagingService(httpClient, testApiUrlFmt, testToken)

	message := &SendMessageRequest{}
	result, err := service.sendMessage(message)

	assert.True(t, result.Successful)
	assert.Nil(t, err)
}

func TestMessagingService_ReturnsErrorIfResponseCodeIsNot200(t *testing.T) {
	httpClient := &mockHttpClient{}
	response := sendMessageResponse{Ok: false}
	responseJson, _ := json.Marshal(response)
	responseBody := io.NopCloser(bytes.NewBufferString(string(responseJson)))
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       responseBody,
	}, nil)
	service, _ := newMessagingService(httpClient, testApiUrlFmt, testToken)

	message := &SendMessageRequest{}
	result, err := service.sendMessage(message)

	assert.False(t, result.Successful)
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(result.Description, "unexpected response code: 500"))
}
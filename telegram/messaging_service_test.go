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

const (
	testApiUrlFmt = "https://api.telegram.local/bot%s/%s"
	testToken     = ""
)

func TestSendMessage_ReturnsErrorIfMessageIsNil(t *testing.T) {
	httpClient := &mockHttpClient{}
	service, _ := newMessagingService(httpClient, testApiUrlFmt, testToken)

	result, err := service.sendMessage(nil)

	assert.False(t, result.Successful)
	assert.Error(t, err)
	assert.Equal(t, errNilMessageRequest, err)
}

func TestSendMessage_ActionResultIsSuccessfulIfMessageIsSentSuccessfully(t *testing.T) {
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

func TestSendMessage_ReturnsErrorIfResponseCodeIsNot200(t *testing.T) {
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

func TestSendMessage_ReturnsErrorIfRequestFails(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	service, _ := newMessagingService(httpClient, testApiUrlFmt, testToken)

	message := &SendMessageRequest{}
	result, err := service.sendMessage(message)

	assert.False(t, result.Successful)
	assert.Error(t, err)
	assert.True(t, strings.EqualFold(err.Error(), "http request failed: error"))
}

func TestSendMessage_ReturnsErrorIfUnableToParseResponse(t *testing.T) {
	httpClient := &mockHttpClient{}
	responseBody := io.NopCloser(bytes.NewBufferString(`invalid json`))
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       responseBody,
	}, nil)
	service, _ := newMessagingService(httpClient, testApiUrlFmt, testToken)

	message := &SendMessageRequest{}
	result, err := service.sendMessage(message)

	assert.False(t, result.Successful)
	assert.Error(t, err)
	assert.True(t, strings.EqualFold(result.Description, "failed to unmarshall sendMessage response"))
}

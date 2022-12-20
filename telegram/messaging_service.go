package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type messagingService struct {
	httpClient httpClient
	apiUrlFmt  string
	token      string
}

func newMessagingService(httpClient httpClient, apirUrlFmt, token string) (*messagingService, error) {
	return &messagingService{
		httpClient: httpClient,
		apiUrlFmt:  apirUrlFmt,
		token:      token,
	}, nil
}

func (s *messagingService) sendMessage(message *SendMessageRequest) (*ActionResult, error) {
	result := &ActionResult{
		Successful: false,
	}

	if message == nil {
		result.Description = errNilMessageRequest.Error()
		return result, errNilMessageRequest
	}

	url := fmt.Sprintf(s.apiUrlFmt, s.token, endpointSendMessage)

	bodyJson, err := json.Marshal(message)
	if err != nil {
		result.Description = "failed to marshal send message request"
		return result, errors.Wrap(err, result.Description)
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		result.Description = "failed to create send message request"
		return result, errors.Wrap(err, result.Description)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		result.Description = "http request failed"
		return result, errors.Wrap(err, result.Description)
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		result.Description = fmt.Sprintf("unexpected response code: %d, %s", response.StatusCode, response.Body)
		return result, errors.New(result.Description)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		result.Description = "failed to parse send message response body"
		return result, errors.Wrap(err, result.Description)
	}

	var resp sendMessageResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		result.Description = "failed to unmarshall sendMessage response"
		return result, errors.Wrap(err, result.Description)
	}

	result.Successful = resp.Ok
	result.Description = resp.Description

	return result, nil
}

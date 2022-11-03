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

func (s *messagingService) sendMessage(message *SendMessageRequest) (bool, error) {
	if message == nil {
		return false, errNilMessageRequest
	}

	url := fmt.Sprintf(s.apiUrlFmt, s.token, endpointSendMessage)

	bodyJson, err := json.Marshal(message)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal send request")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return false, errors.Wrap(err, "failed to create send request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "http request failed")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse send response body")
	}

	var resp sendMessageResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall sendMessage response")
	}

	return resp.Ok, nil
}

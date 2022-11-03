package webhook // import "heytobi.dev/fuse/telegram/webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	httpPost              = "POST"
	endpointSetWebhook    = "setWebhook"    // https://core.telegram.org/bots/api#setwebhook
	endpointDeleteWebhook = "deleteWebhook" // https://core.telegram.org/bots/api#deletewebhook
)

var (
	errMissingWebhookUrl = errors.New("a url is required to register a webhook")
)

// NewService ...
func NewService(httpClient httpClient, apiUrlFmt, token string, AllowedUpdates []string) (*Service, error) {
	service := &Service{
		httpClient:     httpClient,
		apiUrlFmt:      apiUrlFmt,
		token:          token,
		AllowedUpdates: AllowedUpdates,
	}

	return service, nil
}

// RegisterWebhook registers the given webhook to listen for updates.
// Returns the result of the request, True on success.
// See https://core.telegram.org/bots/api#setwebhook
func (s *Service) RegisterWebhook(webhook *Webhook) (bool, error) {

	if webhook.Url == "" {
		return false, errMissingWebhookUrl
	}

	url := fmt.Sprintf(s.apiUrlFmt, s.token, endpointSetWebhook)

	if webhook.AllowedUpdates == nil {
		webhook.AllowedUpdates = s.AllowedUpdates
	}

	bodyJson, err := json.Marshal(webhook)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal register webhook request body")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return false, errors.Wrap(err, "failed to create register webhook request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "register webhook http request failed")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse register webhook response body")
	}

	var resp webhookResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall setWebhook response")
	}

	return resp.Ok, nil
}

// DeleteWebhook deletes the registered webhook.
// See https://core.telegram.org/bots/api#deletewebhook
func (s *Service) DeleteWebhook(dropPendingUpdates bool) (bool, error) {
	url := fmt.Sprintf(s.apiUrlFmt, s.token, endpointDeleteWebhook)

	body := deleteWebhookRequest{DropPendingUpdates: dropPendingUpdates}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal delete webhook request body")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return false, errors.Wrap(err, "failed to create register webhook request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "delete webhook http request failed")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse delete webhook response body")
	}

	var resp webhookResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshall delete Webhook response")
	}

	return resp.Ok, nil
}

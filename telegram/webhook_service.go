package telegram // import "heytobi.dev/fuse/telegram"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type webhookService struct {
	httpClient     httpClient
	apiUrlFmt      string
	token          string
	AllowedUpdates []string `json:"allowed_updates"`
}

func newWebhookService(httpClient httpClient, apirUrlFmt, token string, allowedUpdates []string) (*webhookService, error) {
	return &webhookService{
		httpClient:     httpClient,
		apiUrlFmt:      apirUrlFmt,
		token:          token,
		AllowedUpdates: allowedUpdates,
	}, nil
}

func (s *webhookService) registerWebhook(webhook *Webhook) (bool, error) {
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

// deleteWebhook deletes the registered webhook.
// See https://core.telegram.org/bots/api#deletewebhook
func (s *webhookService) deleteWebhook(dropPendingUpdates bool) (bool, error) {
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

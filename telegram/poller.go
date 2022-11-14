package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"heytobi.dev/fuse/job"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Poller is responsible for continuously checking for updates from Telegram using the getUpdates method.
// See https://core.telegram.org/bots/api#getupdates
type Poller struct {
	config      *Config
	httpClient  httpClient
	updatesChan chan *Update
	offset      int
	apiUrlFmt   string
}

func NewPoller(config *Config, httpClient httpClient) (*Poller, error) {
	if config == nil {
		return nil, errNilConfig
	}

	if config.Token == "" {
		return nil, errMissingToken
	}

	return &Poller{
		httpClient:  httpClient,
		config:      config,
		updatesChan: make(chan *Update),
		apiUrlFmt:   deriveBotApiUrlBase(config) + "/bot%s/%s",
	}, nil
}

func (p *Poller) start() error {
	scheduledJob, err := job.NewScheduledJob(p, p.config.PollingIntervalMS)
	if err != nil {
		return errors.Wrap(err, "failed to initialize scheduled job for poller")
	}

	scheduledJob.Start()

	return nil
}

func (p *Poller) getUpdatesChannel() <-chan *Update {
	return p.updatesChan
}

func (p *Poller) getUpdates() ([]*Update, error) {
	url := fmt.Sprintf(p.apiUrlFmt, p.config.Token, endpointGetUpdates)

	requestBody := getUpdatesRequest{
		Offset:         p.offset,
		Limit:          p.config.PollingUpdatesLimit,
		Timeout:        p.config.PollingTimeout,
		AllowedUpdates: p.config.AllowedUpdates,
	}

	bodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct getUpdates request body")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create getUpdates request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := p.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "getUpdates call failed")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse getUpdates response body")
	}

	var updates getUpdatesResponse
	err = json.Unmarshal(responseBody, &updates)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to unmarshal getUpdates response: %s", string(responseBody[:])))
	}

	return updates.Result, nil
}

// Run ...
func (p *Poller) Run() {
	updates, err := p.getUpdates()
	if err != nil {
		logrus.WithError(err).Error("failed to get updates")
	}

	sort.Slice(updates, func(i, j int) bool {
		return updates[i].ID < updates[j].ID
	})

	for _, update := range updates {
		update := update
		go func(u *Update) {
			p.updatesChan <- update
			if update.ID >= p.offset {
				p.offset = update.ID + 1
			}
		}(update)
	}
}

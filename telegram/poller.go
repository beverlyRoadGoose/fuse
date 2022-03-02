package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
)

const (
	defaultCronSchedule = "*/1 * * * *"
)

// A poller is responsible for continuously checking for updates from Telegram using the getUpdates method.
// See https://core.telegram.org/bots/api#getupdates
type Poller struct {
	httpClient     httpClient
	cronSchedule   string
	updatesChan    chan<- Update
	updatesRequest *http.Request
}

func NewPoller(httpClient httpClient, config *Config) (*Poller, error) {
	if config.Token == "" {
		return nil, errMissingToken
	}

	cronSchedule := config.PollingCronSchedule
	if config.PollingCronSchedule == "" {
		cronSchedule = defaultCronSchedule
	}

	url := fmt.Sprintf(telegramApiUrlFmt, config.Token, getUpdates)

	requestBody, err := json.Marshal(getUpdatesRequest{
		Limit:          config.PollingUpdatesLimit,
		Timeout:        config.PollingTimeout,
		AllowedUpdates: config.AllowedUpdates,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct getUpdates request body")
	}

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create getUpdates request")
	}
	request.Header.Set("Content-Type", "application/json")

	return &Poller{
		httpClient:     httpClient,
		cronSchedule:   cronSchedule,
		updatesChan:    make(chan Update),
		updatesRequest: request,
	}, nil
}

func (p *Poller) start() error {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Cron(p.cronSchedule).Do(func() {
		response, err := p.httpClient.Do(p.updatesRequest)
		if err != nil {
			log.Printf("getUpdates call failed: %s", err.Error())
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("failed to parse getUpdates response body: %s", err.Error())
		}

		var updates []Update
		err = json.Unmarshal(responseBody, &updates)
		if err != nil {
			log.Printf("failed to unmarshal getUpdates response: %s", err.Error())
		}

		for _, update := range updates {
			p.updatesChan <- update
		}
	})
	if err != nil {
		return errors.Wrap(err, "failed to schedule update polling")
	}

	scheduler.StartAsync()

	return nil
}

func (p *Poller) getUpdatesChanel() chan<- Update {
	return p.updatesChan
}

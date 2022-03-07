package telegram // import "heytobi.dev/fuse/telegram"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
	config       *Config
	httpClient   httpClient
	cronSchedule string
	updatesChan  chan *Update
	offset       int
}

func NewPoller(config *Config, httpClient httpClient) (*Poller, error) {
	if config == nil {
		return nil, errNilConfig
	}

	if config.Token == "" {
		return nil, errMissingToken
	}

	cronSchedule := config.PollingCronSchedule
	if config.PollingCronSchedule == "" {
		cronSchedule = defaultCronSchedule
	}

	return &Poller{
		httpClient:   httpClient,
		config:       config,
		cronSchedule: cronSchedule,
		updatesChan:  make(chan *Update, config.PollingUpdatesLimit),
	}, nil
}

func (p *Poller) start() error {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Cron(p.cronSchedule).Do(func() {
		updates, err := p.getUpdates()
		if err != nil {
			logrus.WithError(err).Error("failed to get updates")
		}

		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{})
		l.Info(updates)

		for _, update := range updates {
			update := update
			go func(u *Update) {
				p.updatesChan <- update
				l.Info(fmt.Sprintf("sent %d", update.ID))
				if update.ID >= p.offset {
					p.offset = update.ID + 1
				}
			}(update)
		}
	})
	if err != nil {
		return errors.Wrap(err, "failed to schedule update polling")
	}

	scheduler.StartAsync()

	return nil
}

func (p *Poller) getUpdatesChannel() <-chan *Update {
	return p.updatesChan
}

func (p *Poller) getUpdates() ([]*Update, error) {
	url := fmt.Sprintf(telegramApiUrlFmt, p.config.Token, getUpdates)

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

	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.Info("making request")
	l.Info(string(bodyJson[:]))

	request, err := http.NewRequest(httpPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create getUpdates request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := p.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "getUpdates call failed")
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
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

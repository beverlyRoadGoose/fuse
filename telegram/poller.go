package telegram // import "heytobi.dev/fuse/telegram"
import (
	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"time"
)

const (
	defaultCronSchedule = "*/1 * * * *"
)

// A poller is responsible for continuously checking for updates from Telegram using the getUpdates method.
// See https://core.telegram.org/bots/api#getupdates
type poller struct {
	httpClient     httpClient
	timeout        int
	allowedUpdates []string
	cronSchedule   string
}

func newPoller(httpClient httpClient, timeout int, allowedUpdates []string, cronSchedule string) *poller {
	if cronSchedule == "" {
		cronSchedule = defaultCronSchedule
	}

	return &poller{
		httpClient:     httpClient,
		timeout:        timeout,
		allowedUpdates: allowedUpdates,
		cronSchedule:   cronSchedule,
	}
}

func (p *poller) start() error {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Cron(p.cronSchedule).Do(func() {
		// make update call, push results to channel
	})
	if err != nil {
		return errors.Wrap(err, "failed to schedule update polling")
	}

	return nil
}

package job // import "heytobi.dev/fuse/job"
import (
	"github.com/pkg/errors"
	"time"
)

var (
	errInvalidInterval = errors.New("interval must be at least a second")
	errNilJob          = errors.New("job cannot be nil")
)

// Job defines a schedule-able job
type job interface {
	Run()
}

// ScheduledJob defines a job that is run every n milliseconds, n being defined in intervalMS. The job is executed the
// first time when Start is called
type ScheduledJob struct {
	job        job
	intervalMS int64
}

func NewScheduledJob(job job, intervalMS int64) (*ScheduledJob, error) {
	if job == nil {
		return nil, errNilJob
	}

	if intervalMS < 1000 {
		return nil, errInvalidInterval
	}

	return &ScheduledJob{
		job:        job,
		intervalMS: intervalMS,
	}, nil
}

func (j *ScheduledJob) Start() {
	go func() {
		for {
			j.job.Run()
			time.Sleep(time.Duration(j.intervalMS) * time.Millisecond)
		}
	}()
}

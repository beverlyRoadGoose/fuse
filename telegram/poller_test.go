package telegram // import "heytobi.dev/fuse/telegram"
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPoller_ReturnErrorIfConfigIsNil(t *testing.T) {
	poller, err := NewPoller(nil, nil)

	assert.Nil(t, poller)
	assert.Error(t, err)
	assert.Equal(t, errNilConfig, err)
}

func TestNewPoller_ReturnErrorIfTokenIsMissing(t *testing.T) {
	poller, err := NewPoller(&Config{}, nil)

	assert.Nil(t, poller)
	assert.Error(t, err)
	assert.Equal(t, errMissingToken, err)
}

func TestNewPoller_UseDefaultScheduleIfNoneIsSpecified(t *testing.T) {
	poller, err := NewPoller(&Config{Token: "test"}, &mockHttpClient{})

	assert.Nil(t, err)
	assert.Equal(t, defaultCronSchedule, poller.cronSchedule)
}

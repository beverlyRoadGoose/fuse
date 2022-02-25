package bot // import "heytobi.dev/fuse/bot"
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBot_ReturnErrorIfConfigIsNil(t *testing.T) {
	bot, err := NewBot(nil, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingConfig)
}

func TestNewBot_ReturnErrorIfServiceProviderIsMissing(t *testing.T) {
	bot, err := NewBot(&Config{}, nil)

	assert.Nil(t, bot)
	assert.Error(t, err)
	assert.Equal(t, err, errMissingServiceProvider)
}

package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewScheduledJob_ReturnErrorIfJobIsNil(t *testing.T) {
	job, err := NewScheduledJob(nil, 0)

	assert.Nil(t, job)
	assert.Error(t, err)
	assert.Equal(t, errNilJob, err)
}

func TestNewScheduledJob_ReturnErrorIfIntervalIsLessThan1Second(t *testing.T) {
	job, err := NewScheduledJob(&mockJob{}, 0)

	assert.Nil(t, job)
	assert.Error(t, err)
	assert.Equal(t, errInvalidInterval, err)
}

package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGetUpdates_ReturnErrorIfRequestFails(t *testing.T) {
	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(nil, errors.New("fails"))

	poller, _ := NewPoller(&Config{Token: "test"}, httpClient)
	updates, err := poller.getUpdates()

	assert.Nil(t, updates)
	assert.Error(t, err)
}

func TestGetUpdates_GetUpdatesSuccessfully(t *testing.T) {
	response := getUpdatesResponse{
		Ok: false,
		Result: []*Update{
			{
				ID: 1,
			},
			{
				ID: 1,
			},
		},
	}
	responseJson, _ := json.Marshal(response)
	body := io.NopCloser(bytes.NewBuffer(responseJson))

	httpClient := &mockHttpClient{}
	httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{Body: body}, nil)

	poller, _ := NewPoller(&Config{Token: "test"}, httpClient)
	updates, err := poller.getUpdates()

	assert.NotNil(t, updates)
	assert.Equal(t, 1, updates[0].ID)
	assert.NoError(t, err)
}

package telegram

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type mockHttpClient struct {
	mock.Mock
}

func (m *mockHttpClient) Do(request *http.Request) (*http.Response, error) {
	args := m.Called(request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*http.Response), args.Error(1)
}

package bot // import "heytobi.dev/fuse/bot"

import (
	"github.com/stretchr/testify/mock"
)

type mockMessagingServiceProvider struct {
	mock.Mock
}

func (m *mockMessagingServiceProvider) Start() error {
	args := m.Called()

	return args.Error(0)
}

func (m *mockMessagingServiceProvider) Send(message Sendable) error {
	args := m.Called(message)

	return args.Error(0)
}

func (m *mockMessagingServiceProvider) RegisterHandler(command string, handlerFunc HandlerFunc) error {
	args := m.Called(command, handlerFunc)

	return args.Error(0)
}

func (m *mockMessagingServiceProvider) ProcessUpdate(update Update) error {
	args := m.Called(update)

	return args.Error(0)
}

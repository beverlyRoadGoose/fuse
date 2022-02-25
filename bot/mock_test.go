package bot // import "heytobi.dev/fuse/bot"

import "github.com/stretchr/testify/mock"

type mockMessagingServiceProvider struct {
	mock.Mock
}

func (m *mockMessagingServiceProvider) SendMessage() error {
	args := m.Called()

	return args.Error(0)
}

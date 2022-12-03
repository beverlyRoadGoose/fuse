package conversation

import (
	"github.com/stretchr/testify/mock"
	"heytobi.dev/fuse/telegram"
)

type mockSequence struct {
	mock.Mock
}

type mockBot struct {
	mock.Mock
}

func (m *mockSequence) Start(orchestrator Orchestrator) {
	m.Called(orchestrator)
}

func (m *mockSequence) Process(update *telegram.Update) error {
	args := m.Called(update)
	return args.Error(0)
}

func (m *mockSequence) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockBot) SendMessage(message *telegram.SendMessageRequest) (*telegram.ActionResult, error) {
	args := m.Called(message)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*telegram.ActionResult), args.Error(1)
}

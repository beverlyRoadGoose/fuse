package conversation

import (
	"github.com/stretchr/testify/mock"
	"heytobi.dev/fuse/telegram"
)

type mockSequence struct {
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

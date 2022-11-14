package job

import "github.com/stretchr/testify/mock"

type mockJob struct {
	mock.Mock
}

func (m *mockJob) Run() {
}

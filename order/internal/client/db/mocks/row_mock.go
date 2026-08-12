package mocks

import (
	"github.com/stretchr/testify/mock"
)

type RowMock struct {
	mock.Mock
}

func (m *RowMock) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

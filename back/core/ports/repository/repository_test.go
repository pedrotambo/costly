package repository_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

type clockMock struct {
	mock.Mock
}

func (m *clockMock) Now() time.Time {
	args := m.Called()
	value := args.Get(0)
	now, ok := value.(time.Time)
	if !ok {
		panic(fmt.Errorf("Error getting now"))
	}
	return now
}

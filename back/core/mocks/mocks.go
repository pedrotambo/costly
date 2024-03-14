package mocks

import (
	"context"
	"costly/core/ports/database"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

type ClockMock struct {
	mock.Mock
}

func (m *ClockMock) Now() time.Time {
	args := m.Called()
	value := args.Get(0)
	now, ok := value.(time.Time)
	if !ok {
		panic(fmt.Errorf("Error getting now"))
	}
	return now
}

var ErrDBInternal = errors.New("internal db error")

type DatabaseMock struct {
	mock.Mock
}

type errorRow struct{}

func (e *errorRow) Scan(dest ...any) error {
	return ErrDBInternal
}

func (e *errorRow) Next() bool {
	return false
}

func (dm *DatabaseMock) QueryRowContext(ctx context.Context, query string, args ...any) database.RowScanner {
	return &errorRow{}
}

func (dm *DatabaseMock) QueryContext(ctx context.Context, query string, args ...any) (database.RowsScanner, error) {
	return nil, ErrDBInternal
}

func (dm *DatabaseMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, ErrDBInternal
}

func (dm *DatabaseMock) WithTx(ctx context.Context, op func(tx database.TX) error) error {
	return ErrDBInternal
}

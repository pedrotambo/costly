package database

import (
	"context"
	"database/sql"
)

type RowScanner interface {
	Scan(dest ...any) error
}

type RowsScanner interface {
	RowScanner
	Next() bool
}

type dbSession interface {
	QueryRowContext(ctx context.Context, query string, args ...any) RowScanner
	QueryContext(ctx context.Context, query string, args ...any) (RowsScanner, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type sqlDBSession interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type databaseSession struct {
	session sqlDBSession
}

func newDBSession(db sqlDBSession) dbSession {
	return &databaseSession{session: db}
}

func (db *databaseSession) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	return db.session.QueryRowContext(ctx, query, args...)
}

func (db *databaseSession) QueryContext(ctx context.Context, query string, args ...any) (RowsScanner, error) {
	return db.session.QueryContext(ctx, query, args...)
}

func (db *databaseSession) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.session.ExecContext(ctx, query, args...)
}

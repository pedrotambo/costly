package database

import (
	"context"
	"costly/core/ports/logger"
	sql2 "costly/sql"
	"database/sql"
	"fmt"
)

type RowScanner interface {
	Scan(dest ...any) error
}

type RowsScanner interface {
	RowScanner
	Next() bool
}

type RowQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) RowScanner
}

type RowsQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (RowsScanner, error)
}

type QueryExecuter interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Database interface {
	TX
	WithTx(ctx context.Context, op func(tx TX) error) error
}

type TX interface {
	RowQuerier
	RowsQuerier
	QueryExecuter
}

type database struct {
	sqlDB  *sql.DB
	logger logger.Logger
}

func New(connectionString string, logger logger.Logger) (Database, error) {
	return NewFromDatasource(fmt.Sprintf("file:%s?_foreign_keys=on", connectionString), logger)
}

func NewFromDatasource(datasourceName string, logger logger.Logger) (Database, error) {
	db, err := sql.Open("sqlite3", datasourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to db: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	_, err = sql2.RunMigrations(db, logger)
	if err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	return &database{db, logger}, nil
}

func (db *database) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	return db.sqlDB.QueryRowContext(ctx, query, args...)
}

func (db *database) QueryContext(ctx context.Context, query string, args ...any) (RowsScanner, error) {
	return db.sqlDB.QueryContext(ctx, query, args...)
}

func (db *database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.sqlDB.ExecContext(ctx, query, args...)
}

func (db *database) WithTx(ctx context.Context, op func(tx TX) error) error {
	tx, err := db.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	err = op(NewTX(tx))
	if err != nil {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			db.logger.Error(err, "failed to rollback transaction")
		}
		return err
	}

	return tx.Commit()
}

type dbtx struct {
	sqlTx *sql.Tx
}

func NewTX(sqltx *sql.Tx) TX {
	return &dbtx{sqlTx: sqltx}
}

func (tx *dbtx) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	return tx.sqlTx.QueryRowContext(ctx, query, args...)
}

func (tx *dbtx) QueryContext(ctx context.Context, query string, args ...any) (RowsScanner, error) {
	return tx.sqlTx.QueryContext(ctx, query, args...)
}

func (tx *dbtx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.sqlTx.ExecContext(ctx, query, args...)
}

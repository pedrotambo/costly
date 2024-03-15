package database

import (
	"context"
	"costly/core/ports/logger"
	sql2 "costly/sql"
	"database/sql"
	"fmt"
)

type Database interface {
	dbSession
	WithTx(ctx context.Context, op func(tx Database) error) error
}

func New(connectionString string, logger logger.Logger) (Database, error) {
	return NewFromDatasource(fmt.Sprintf("file:%s", connectionString), logger)
}

func NewFromDatasource(datasourceName string, logger logger.Logger) (Database, error) {
	db, err := sql.Open("sqlite3", datasourceName+"?_foreign_keys=on")
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

	return newPoolDB(db), nil
}

type pooldb struct {
	sqlDB *sql.DB
	dbSession
}

func newPoolDB(sqlDB *sql.DB) *pooldb {
	return &pooldb{
		sqlDB:     sqlDB,
		dbSession: newDBSession(sqlDB),
	}
}

func (db *pooldb) WithTx(ctx context.Context, op func(tx Database) error) error {
	tx, err := db.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	err = op(newTX(tx))
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type dbtx struct {
	sqlTx *sql.Tx
	dbSession
}

func newTX(sqltx *sql.Tx) Database {
	return &dbtx{
		sqlTx:     sqltx,
		dbSession: newDBSession(sqltx),
	}
}

func (tx *dbtx) WithTx(ctx context.Context, op func(tx Database) error) error {
	return op(tx)
}

package database

import (
	"context"
	sql2 "costly/sql"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

type Database struct {
	*sql.DB
}

func New(connectionString string) (*Database, error) {
	return NewFromDatasource(fmt.Sprintf("file:%s?_foreign_keys=on", connectionString))
}

func NewFromDatasource(datasourceName string) (*Database, error) {
	db, err := sql.Open("sqlite3", datasourceName)

	if err != nil {
		return nil, fmt.Errorf("failed to open connection to db: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	_, err = sql2.RunMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	return &Database{db}, nil
}

func (database *Database) WithTx(ctx context.Context, op func(tx *sql.Tx) error) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	err = op(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to rollback transaction")
		}
		return err
	}

	return tx.Commit()
}

package sql

import (
	"costly/core/components/logger"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

// RunMigrations receives a *sql.DB instance with a MySQL backend, and runs the appropriated migrations.
func RunMigrations(db *sql.DB, logger logger.Logger) (string, error) {
	logger.Info("running migrations")

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return "", fmt.Errorf("error creating migrations driver: %w", err)
	}

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return "", fmt.Errorf("error loading migration files: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", d, "mysql", driver)
	if err != nil {
		return "", fmt.Errorf("error creating migrator: %w", err)
	}

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return "", fmt.Errorf("error executing migrations: %w", err)
	}

	version, _, _ := migrator.Version()

	logger.Info("migrations run")
	return fmt.Sprint(version), nil
}

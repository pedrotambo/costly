package ports

import (
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"fmt"
)

type Ports struct {
	Database database.Database
	Clock    clock.Clock
	Logger   logger.Logger
}

func New(logLevel string, connectionString string) (*Ports, error) {
	logger, err := logger.New(logLevel)
	if err != nil {
		return &Ports{}, fmt.Errorf("could not create logger, Err: %s", err)
	}
	database, err := database.New(connectionString, logger)
	if err != nil {
		return &Ports{}, fmt.Errorf("could not initialize database. Err: %s", err)
	}
	clock := clock.New()
	return &Ports{
		Database: database,
		Clock:    clock,
		Logger:   logger,
	}, nil
}

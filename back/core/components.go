package core

import (
	"costly/api"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"fmt"
	"net/http"
	"os"
)

type AppComponents struct {
	Logger     logger.Logger
	Database   database.Database
	Server     *http.Server
	Clock      clock.Clock
	Repository rpst.Repository
}

func InitComponents(config *Config) (*AppComponents, error) {
	logger, err := logger.New(config.LogLevel)
	if err != nil {
		fmt.Printf("Could not create logger. Err: %s\n", err)
		os.Exit(1)
	}

	logger.Info("Running server...")

	database, err := database.New(config.Database.ConnectionString, logger)
	if err != nil {
		logger.Error(err, "could not initialize database")
		os.Exit(1)
	}

	loggerMiddleware := api.NewLoggerMiddleware(logger)
	authMiddleware := api.NewAuthMiddleware([]byte(config.AuthSecret), logger)

	clock := clock.New()
	repository := rpst.New(database, clock, logger)
	useCases := usecases.New(repository, clock)
	router := api.NewRouter(repository, useCases, authMiddleware, loggerMiddleware)
	server := http.Server{
		Addr:    config.ListenAddress,
		Handler: router,
	}

	return &AppComponents{
		Logger:     logger,
		Database:   database,
		Server:     &server,
		Clock:      clock,
		Repository: repository,
	}, nil
}

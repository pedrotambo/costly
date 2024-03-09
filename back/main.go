package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"costly/api"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Could not load configuration. Err: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(*config)

	components, err := initComponents(config)

	if err != nil {
		fmt.Printf("Could not initialize components. Err: %s\n", err)
		os.Exit(1)
	}

	done := make(chan bool, 1)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel

		if err := components.server.Shutdown(context.Background()); err != nil {
			components.logger.Error(err, "could not gracefully shutdown server")
		}

		done <- true
	}()

	if err := components.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		components.logger.Error(err, "could not start server")
		os.Exit(1)
	}

	<-done

	components.logger.Info("app stopped")

	fmt.Println(done)
}

type AppComponents struct {
	logger     logger.Logger
	database   database.Database
	server     *http.Server
	clock      clock.Clock
	repository rpst.Repository
}

func initComponents(config *Config) (AppComponents, error) {
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
	router := api.NewRouter(repository, authMiddleware, loggerMiddleware)
	server := http.Server{
		Addr:    config.ListenAddress,
		Handler: router,
	}

	return AppComponents{
		logger:     logger,
		database:   database,
		server:     &server,
		clock:      clock,
		repository: repository,
	}, nil
}

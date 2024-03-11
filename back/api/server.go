package api

import (
	"context"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Start()
}

type server struct {
	httpServer   *http.Server
	portAdapters *usecases.Ports
}

func NewServer(config *Config) Server {
	portAdapters, err := initPortAdapters(config)
	if err != nil {
		fmt.Printf("Could not initialize components. Err: %s\n", err)
		os.Exit(1)
	}
	logger := portAdapters.Logger
	loggerMiddleware := NewLoggerMiddleware(logger)
	authMiddleware := NewAuthMiddleware([]byte(config.AuthSecret), logger)
	router := NewRouter(usecases.New(portAdapters), authMiddleware, loggerMiddleware)
	return &server{
		httpServer: &http.Server{
			Addr:    config.ListenAddress,
			Handler: router,
		},
		portAdapters: portAdapters,
	}
}

func (s *server) Start() {
	done := make(chan bool, 1)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	logger := s.portAdapters.Logger
	logger.Info("Running server...")

	go func() {
		<-signalChannel

		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			logger.Error(err, "could not gracefully shutdown server")
		}

		done <- true
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err, "could not start server")
		os.Exit(1)
	}

	<-done

	logger.Info("app stopped")

	fmt.Println(done)
}

func initPortAdapters(config *Config) (*usecases.Ports, error) {
	logger, err := logger.New(config.LogLevel)
	if err != nil {
		fmt.Printf("Could not create logger. Err: %s\n", err)
		os.Exit(1)
	}
	database, err := database.New(config.Database.ConnectionString, logger)
	if err != nil {
		logger.Error(err, "could not initialize database")
		os.Exit(1)
	}
	clock := clock.New()
	repository := rpst.New(database, clock, logger)
	return &usecases.Ports{
		Logger:     logger,
		Repository: repository,
		Clock:      clock,
	}, nil
}

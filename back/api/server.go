package api

import (
	"context"
	comps "costly/core/components"
	"costly/core/components/logger"
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
	httpServer *http.Server
	logger     logger.Logger
}

func NewServer(config *Config) Server {
	components, err := comps.InitComponents(&config.ComponentsConfig)
	logger := components.Logger
	if err != nil {
		fmt.Printf("Could not initialize components. Err: %s\n", err)
		os.Exit(1)
	}
	loggerMiddleware := NewLoggerMiddleware(logger)
	authMiddleware := NewAuthMiddleware([]byte(config.AuthSecret), logger)
	router := NewRouter(components, authMiddleware, loggerMiddleware)
	return &server{
		httpServer: &http.Server{
			Addr:    config.ListenAddress,
			Handler: router,
		},
		logger: logger,
	}
}

func (s *server) Start() {
	done := make(chan bool, 1)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	logger := s.logger
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

package api

import (
	"context"
	comps "costly/core/components"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/components/recipes"
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
	components, err := initComponents(config)
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

func initComponents(config *Config) (*comps.Components, error) {
	logger, err := logger.New(config.LogLevel)
	if err != nil {
		fmt.Printf("Could not create logger. Err: %s\n", err)
		return &comps.Components{}, err
	}
	database, err := database.New(config.Database.ConnectionString, logger)
	if err != nil {
		logger.Error(err, "could not initialize database")
		return &comps.Components{}, err
	}
	clock := clock.New()
	ingredientComponent := ingredients.New(database, clock, logger)
	recipeComponent := recipes.New(database, clock, logger, ingredientComponent)
	return &comps.Components{
		IngredientComponent: ingredientComponent,
		RecipeComponent:     recipeComponent,
		Logger:              logger,
		Database:            database,
		Clock:               clock,
	}, nil
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"costly/api"
	costs "costly/core/logic"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/repository"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Could not load configuration. Err: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(config)

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
			components.logger.Error().Err(err).Msg("could not gracefully shutdown server")
		}

		done <- true
	}()

	if err := components.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		components.logger.Error().Err(err).Msg("could not start server")
		os.Exit(1)
	}

	<-done

	components.logger.Info().Msg("app stopped")

	fmt.Println(done)
}

type AppComponents struct {
	logger      *zerolog.Logger
	database    *database.Database
	router      chi.Router
	server      *http.Server
	clock       clock.Clock
	repository  *repository.Repository
	costService costs.CostService
}

func initComponents(config *Config) (AppComponents, error) {
	zLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		fmt.Printf("Could not parse log level. Err: %s\n", err)
		os.Exit(1)
	}
	logger := zerolog.New(os.Stderr).Level(zLevel).With().Timestamp().Logger()

	logger.Info().Msg("Running server...")

	database, err := database.New(config.Database.ConnectionString)
	if err != nil {
		logger.Error().Err(err).Msg("could not initialize database")
		os.Exit(1)
	}

	// logger.Info().Str("version", build.Version).Msg("app started")

	loggerInjectorMiddleware := api.Middleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithContext(r.Context())
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	authSupport := api.NewAuthSupport([]byte(config.AuthSecret))
	// TODO: Add some flag to enable just in dev
	authSupport.PrintDebug(logger)

	clock := clock.New()
	repository := repository.New(database, clock)
	costService := costs.NewCostService()
	router := api.NewRouter(database, clock, repository, costService, authSupport, loggerInjectorMiddleware)
	server := http.Server{
		Addr:    config.ListenAddress,
		Handler: router,
	}

	return AppComponents{
		logger:      &logger,
		database:    database,
		router:      router,
		server:      &server,
		clock:       clock,
		repository:  repository,
		costService: costService,
	}, nil
}

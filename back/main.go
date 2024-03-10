package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"costly/core"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := core.LoadConfig()
	if err != nil {
		fmt.Printf("Could not load configuration. Err: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(*config)

	components, err := core.InitComponents(config)

	if err != nil {
		fmt.Printf("Could not initialize components. Err: %s\n", err)
		os.Exit(1)
	}

	done := make(chan bool, 1)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel

		if err := components.Server.Shutdown(context.Background()); err != nil {
			components.Logger.Error(err, "could not gracefully shutdown server")
		}

		done <- true
	}()

	if err := components.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		components.Logger.Error(err, "could not start server")
		os.Exit(1)
	}

	<-done

	components.Logger.Info("app stopped")

	fmt.Println(done)
}

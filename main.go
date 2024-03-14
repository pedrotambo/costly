package main

import (
	"fmt"
	"os"

	"costly/api"
	"costly/core/ports"
	comps "costly/core/usecases"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Could not load configuration. Err: %s\n", err)
		os.Exit(1)
	}
	ports, err := ports.New(config.LogLevel, config.Database.ConnectionString)
	if err != nil {
		fmt.Printf("Could not initialize adapters. Err: %s\n", err)
		os.Exit(1)
	}
	components, err := comps.New(ports)
	if err != nil {
		fmt.Printf("Could not initialize components. Err: %s\n", err)
		os.Exit(1)
	}
	api.NewServer(config.ListenAddress, config.AuthSecret, components, ports.Logger).Start()
}

package main

import (
	"fmt"
	"os"

	"costly/api"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := api.LoadConfig()
	if err != nil {
		fmt.Printf("Could not load configuration. Err: %s\n", err)
		os.Exit(1)
	}
	api.NewServer(config).Start()
}

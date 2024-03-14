package main

import (
	"flag"
	"fmt"
)

type Database struct {
	ConnectionString string
}

type Config struct {
	LogLevel      string
	ListenAddress string
	AuthSecret    string
	Database      struct {
		ConnectionString string
	}
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	fs := flag.CommandLine
	fs.StringVar(&cfg.ListenAddress, "listen-addr", ":3000", "Main listen address for the HTTP server.")
	fs.StringVar(&cfg.AuthSecret, "auth-secret", "sample-secret", "Authentication secret for signing JWTs.")
	fs.StringVar(&cfg.Database.ConnectionString, "db.connection-string", "", "SQLite connection string.")
	fs.StringVar(&cfg.LogLevel, "log.level", "info", "Log level.")
	flag.Parse()
	if cfg.Database.ConnectionString == "" {
		return &Config{}, fmt.Errorf("empty DB connection string")
	}
	fmt.Println(cfg)
	return &cfg, nil
}

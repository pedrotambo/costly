package main

import (
	"flag"
	"fmt"
)

type Config struct {
	LogLevel      string
	ListenAddress string
	AuthSecret    string
	Database      struct {
		ConnectionString string
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.registerFlags(flag.CommandLine)
	flag.Parse()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) registerFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.LogLevel, "log.level", "info", "Log level.")
	fs.StringVar(&cfg.ListenAddress, "listen-addr", ":3000", "Main listen address for the HTTP server.")
	fs.StringVar(&cfg.Database.ConnectionString, "db.connection-string", "", "SQLite connection string.")
	fs.StringVar(&cfg.AuthSecret, "auth-secret", "sample-secret", "Authentication secret for signing JWTs.")
}

func (cfg *Config) Validate() error {
	if cfg.Database.ConnectionString == "" {
		return fmt.Errorf("empty DB connection string")
	}
	return nil
}

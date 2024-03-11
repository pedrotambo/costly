package api

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
	cfg := Config{}
	fs := flag.CommandLine
	fs.StringVar(&cfg.LogLevel, "log.level", "info", "Log level.")
	fs.StringVar(&cfg.ListenAddress, "listen-addr", ":3000", "Main listen address for the HTTP server.")
	fs.StringVar(&cfg.Database.ConnectionString, "db.connection-string", "", "SQLite connection string.")
	fs.StringVar(&cfg.AuthSecret, "auth-secret", "sample-secret", "Authentication secret for signing JWTs.")
	flag.Parse()
	if cfg.Database.ConnectionString == "" {
		return &Config{}, fmt.Errorf("empty DB connection string")
	}
	fmt.Println(cfg)
	return &cfg, nil
}

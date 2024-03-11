package api

import (
	comps "costly/core/components"
	"flag"
	"fmt"
)

type Config struct {
	ListenAddress    string
	AuthSecret       string
	ComponentsConfig comps.Config
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	fs := flag.CommandLine
	fs.StringVar(&cfg.ListenAddress, "listen-addr", ":3000", "Main listen address for the HTTP server.")
	fs.StringVar(&cfg.AuthSecret, "auth-secret", "sample-secret", "Authentication secret for signing JWTs.")
	fs.StringVar(&cfg.ComponentsConfig.Database.ConnectionString, "db.connection-string", "", "SQLite connection string.")
	fs.StringVar(&cfg.ComponentsConfig.LogLevel, "log.level", "info", "Log level.")
	flag.Parse()
	if cfg.ComponentsConfig.Database.ConnectionString == "" {
		return &Config{}, fmt.Errorf("empty DB connection string")
	}
	fmt.Println(cfg)
	return &cfg, nil
}

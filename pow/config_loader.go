package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:  8080,
		DebugMode:   false,
		DatabaseURL: "localhost:5432",
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.ServerPort = port
		}
	}

	if debugStr := os.Getenv("DEBUG_MODE"); debugStr != "" {
		cfg.DebugMode = strings.ToLower(debugStr) == "true"
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}

	if hosts := os.Getenv("ALLOWED_HOSTS"); hosts != "" {
		cfg.AllowedHosts = strings.Split(hosts, ",")
	}

	return cfg, nil
}
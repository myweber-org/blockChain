package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	LogLevel   string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort: 8080,
		DebugMode:  false,
		LogLevel:   "info",
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.ServerPort = port
		}
	}

	if debugStr := os.Getenv("DEBUG_MODE"); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			cfg.DebugMode = debug
		}
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return cfg, nil
}
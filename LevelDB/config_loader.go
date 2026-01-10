package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	FeatureFlags map[string]bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		FeatureFlags: make(map[string]bool),
	}

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port

	cfg.DBHost = getEnv("DB_HOST", "localhost")

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort

	debug, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug

	flags := strings.Split(getEnv("FEATURE_FLAGS", ""), ",")
	for _, flag := range flags {
		trimmed := strings.TrimSpace(flag)
		if trimmed != "" {
			cfg.FeatureFlags[trimmed] = true
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
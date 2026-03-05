
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	APIKeys    []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port
	
	debugStr := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	
	apiKeysStr := os.Getenv("API_KEYS")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}
	
	return cfg, nil
}

func ValidateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if cfg.DatabaseURL == "" {
		return errors.New("database URL cannot be empty")
	}
	
	return nil
}
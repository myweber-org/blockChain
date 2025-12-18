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
	APIKeys    []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.ServerPort = getEnvAsInt("SERVER_PORT", 8080)
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnvAsInt("DB_PORT", 5432)
	cfg.DebugMode = getEnvAsBool("DEBUG_MODE", false)
	cfg.APIKeys = getEnvAsSlice("API_KEYS", []string{}, ",")

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
	}
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
	}
	if len(cfg.APIKeys) == 0 {
		return &ConfigError{Field: "API_KEYS", Message: "at least one API key is required"}
	}
	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: field " + e.Field + " - " + e.Message
}
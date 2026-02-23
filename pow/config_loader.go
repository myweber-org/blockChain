package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort   int
    DatabaseURL  string
    LogLevel     string
    CacheEnabled bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:     getEnv("LOG_LEVEL", "info"),
        CacheEnabled: getEnvAsBool("CACHE_ENABLED", true),
    }

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

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }

    if !strings.HasPrefix(cfg.DatabaseURL, "postgres://") {
        return &ConfigError{Field: "DatabaseURL", Message: "must be a PostgreSQL connection string"}
    }

    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return &ConfigError{Field: "LogLevel", Message: "must be one of: debug, info, warn, error"}
    }

    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}
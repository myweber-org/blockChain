package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort int
    DatabaseURL string
    LogLevel string
    EnableCache bool
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        EnableCache: getEnvAsBool("ENABLE_CACHE", true),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
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
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    strValue = strings.ToLower(strValue)
    return strValue == "true" || strValue == "1" || strValue == "yes"
}

func validateConfig(cfg *AppConfig) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.ServerPort)
    }
    if cfg.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }
    allowedLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    if !allowedLogLevels[strings.ToLower(cfg.LogLevel)] {
        return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
    }
    return nil
}
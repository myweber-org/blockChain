package config

import (
    "fmt"
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

    port, err := getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    cfg.ServerPort = port

    cfg.DBHost = getStringEnv("DB_HOST", "localhost")

    dbPort, err := getIntEnv("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }
    cfg.DBPort = dbPort

    debug, err := getBoolEnv("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }
    cfg.DebugMode = debug

    cfg.APIKeys = getStringSliceEnv("API_KEYS", []string{})

    return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, err
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return false, err
        }
        return boolValue, nil
    }
    return defaultValue, nil
}

func getStringSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}
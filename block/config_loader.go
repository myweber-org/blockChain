package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        CacheEnabled:   true,
        MaxConnections: 100,
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
            cfg.ServerPort = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheFlag := os.Getenv("CACHE_ENABLED"); cacheFlag != "" {
        cfg.CacheEnabled = strings.ToLower(cacheFlag) == "true"
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        if maxConn, err := strconv.Atoi(maxConnStr); err == nil && maxConn > 0 {
            cfg.MaxConnections = maxConn
        }
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Value: cfg.ServerPort}
    }
    if cfg.DatabaseURL == "" {
        return &ConfigError{Field: "DatabaseURL", Value: "empty"}
    }
    if cfg.MaxConnections < 1 {
        return &ConfigError{Field: "MaxConnections", Value: cfg.MaxConnections}
    }
    return nil
}

type ConfigError struct {
    Field string
    Value interface{}
}

func (e *ConfigError) Error() string {
    return "invalid configuration value for " + e.Field
}
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
    AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        CacheEnabled:   true,
        MaxConnections: 100,
        AllowedOrigins: []string{"http://localhost:3000"},
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
            cfg.ServerPort = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheStr := os.Getenv("CACHE_ENABLED"); cacheStr != "" {
        cfg.CacheEnabled = strings.ToLower(cacheStr) == "true"
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        if maxConn, err := strconv.Atoi(maxConnStr); err == nil && maxConn > 0 {
            cfg.MaxConnections = maxConn
        }
    }

    if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
        cfg.AllowedOrigins = strings.Split(origins, ",")
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return ErrInvalidPort
    }

    if cfg.DatabaseURL == "" {
        return ErrMissingDatabaseURL
    }

    if cfg.MaxConnections < 1 {
        return ErrInvalidMaxConnections
    }

    return nil
}

var (
    ErrInvalidPort           = errors.New("invalid server port")
    ErrMissingDatabaseURL    = errors.New("database URL is required")
    ErrInvalidMaxConnections = errors.New("max connections must be positive")
)
package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableCache bool
    MaxConnections int
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    enableCache := getEnv("ENABLE_CACHE", "true")
    cfg.EnableCache = strings.ToLower(enableCache) == "true"
    
    maxConnStr := getEnv("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
    }
    cfg.MaxConnections = maxConn
    
    origins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
    cfg.AllowedOrigins = strings.Split(origins, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
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

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    
    if cfg.MaxConnections < 1 {
        return fmt.Errorf("max connections must be positive")
    }
    
    if cfg.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }
    
    return nil
}
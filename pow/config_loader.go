package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DebugMode  bool
    DatabaseURL string
    AllowedHosts []string
}

func Load() (*Config, error) {
    cfg := &Config{
        ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
        DebugMode:   getEnvAsBool("DEBUG_MODE", false),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost"}, ","),
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
    if val, err := strconv.ParseBool(valueStr); err == nil {
        return val
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
    CacheEnabled bool
    LogLevel string
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    cacheEnabledStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabledStr) == "true"
    
    cfg.LogLevel = getEnvWithDefault("LOG_LEVEL", "info")
    if !isValidLogLevel(cfg.LogLevel) {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
    }
    
    maxConnStr := getEnvWithDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil || maxConn <= 0 {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS value: %s", maxConnStr)
    }
    cfg.MaxConnections = maxConn
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func isValidLogLevel(level string) bool {
    validLevels := map[string]bool{
        "debug": true,
        "info": true,
        "warn": true,
        "error": true,
        "fatal": true,
    }
    return validLevels[strings.ToLower(level)]
}
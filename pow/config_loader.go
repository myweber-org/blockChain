package config

import (
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            Username: getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            Database: getEnv("DB_NAME", "appdb"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
        Server: ServerConfig{
            Port:         getEnvAsInt("SERVER_PORT", 8080),
            ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
            WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
            DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        },
        LogLevel: getEnv("LOG_LEVEL", "info"),
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
    var errors []string

    if cfg.Database.Host == "" {
        errors = append(errors, "database host cannot be empty")
    }
    if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
        errors = append(errors, "database port must be between 1 and 65535")
    }
    if cfg.Database.Username == "" {
        errors = append(errors, "database username cannot be empty")
    }
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        errors = append(errors, "server port must be between 1 and 65535")
    }
    if cfg.Server.ReadTimeout < 1 {
        errors = append(errors, "read timeout must be positive")
    }
    if cfg.Server.WriteTimeout < 1 {
        errors = append(errors, "write timeout must be positive")
    }

    if len(errors) > 0 {
        return &ConfigError{Errors: errors}
    }
    return nil
}

type ConfigError struct {
    Errors []string
}

func (e *ConfigError) Error() string {
    return "configuration validation failed: " + strings.Join(e.Errors, "; ")
}
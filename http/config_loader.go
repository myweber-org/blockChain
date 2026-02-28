package config

import (
    "fmt"
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
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Debug    bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    dbHost := getEnv("DB_HOST", "localhost")
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %v", err)
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASS", ""),
        Database: getEnv("DB_NAME", "appdb"),
    }

    serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }

    readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid READ_TIMEOUT: %v", err)
    }

    writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %v", err)
    }

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    debugStr := strings.ToLower(getEnv("DEBUG", "false"))
    cfg.Debug = debugStr == "true" || debugStr == "1"

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("database password is required")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return nil, fmt.Errorf("server port must be between 1 and 65535")
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
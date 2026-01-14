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
    dbPort := getEnvAsInt("DB_PORT", 5432)
    dbUser := getEnv("DB_USER", "postgres")
    dbPass := getEnv("DB_PASS", "")
    dbName := getEnv("DB_NAME", "appdb")

    if dbPort < 1 || dbPort > 65535 {
        return nil, fmt.Errorf("invalid database port: %d", dbPort)
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: dbUser,
        Password: dbPass,
        Database: dbName,
    }

    serverPort := getEnvAsInt("SERVER_PORT", 8080)
    readTimeout := getEnvAsInt("READ_TIMEOUT", 30)
    writeTimeout := getEnvAsInt("WRITE_TIMEOUT", 30)

    if serverPort < 1 || serverPort > 65535 {
        return nil, fmt.Errorf("invalid server port: %d", serverPort)
    }

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    debugStr := strings.ToLower(getEnv("DEBUG", "false"))
    cfg.Debug = debugStr == "true" || debugStr == "1"

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

    intValue, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }

    return intValue
}

func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if c.Database.Username == "" {
        return fmt.Errorf("database username cannot be empty")
    }
    if c.Database.Database == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    return nil
}
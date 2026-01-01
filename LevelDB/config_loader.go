package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASS"`
	Database string `yaml:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideString(&config.Database.Host, "DB_HOST")
	overrideInt(&config.Database.Port, "DB_PORT")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")

	overrideInt(&config.Server.Port, "SERVER_PORT")
	overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
	overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
	overrideBool(&config.Server.DebugMode, "DEBUG_MODE")
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}

func overrideBool(field *bool, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val == "true" || val == "1" || val == "yes"
	}
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("invalid database port")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("invalid server port")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("write timeout cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}
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
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type AppConfig struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*AppConfig, error) {
    dbConfig, err := loadDatabaseConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load database config: %w", err)
    }

    serverConfig, err := loadServerConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load server config: %w", err)
    }

    logLevel := getEnvWithDefault("LOG_LEVEL", "info")

    return &AppConfig{
        Database: dbConfig,
        Server:   serverConfig,
        LogLevel: strings.ToLower(logLevel),
    }, nil
}

func loadDatabaseConfig() (DatabaseConfig, error) {
    host := getEnvRequired("DB_HOST")
    portStr := getEnvRequired("DB_PORT")
    username := getEnvRequired("DB_USERNAME")
    password := getEnvRequired("DB_PASSWORD")
    database := getEnvRequired("DB_DATABASE")
    sslMode := getEnvWithDefault("DB_SSL_MODE", "require")

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return DatabaseConfig{}, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    if port < 1 || port > 65535 {
        return DatabaseConfig{}, fmt.Errorf("DB_PORT out of range: %d", port)
    }

    return DatabaseConfig{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        Database: database,
        SSLMode:  sslMode,
    }, nil
}

func loadServerConfig() (ServerConfig, error) {
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    readTimeoutStr := getEnvWithDefault("SERVER_READ_TIMEOUT", "30")
    writeTimeoutStr := getEnvWithDefault("SERVER_WRITE_TIMEOUT", "30")
    debugModeStr := getEnvWithDefault("DEBUG_MODE", "false")

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return ServerConfig{}, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }

    readTimeout, err := strconv.Atoi(readTimeoutStr)
    if err != nil {
        return ServerConfig{}, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
    }

    writeTimeout, err := strconv.Atoi(writeTimeoutStr)
    if err != nil {
        return ServerConfig{}, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
    }

    debugMode, err := strconv.ParseBool(debugModeStr)
    if err != nil {
        return ServerConfig{}, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }

    if port < 1 || port > 65535 {
        return ServerConfig{}, fmt.Errorf("SERVER_PORT out of range: %d", port)
    }

    return ServerConfig{
        Port:         port,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
        DebugMode:    debugMode,
    }, nil
}

func getEnvRequired(key string) string {
    value := os.Getenv(key)
    if value == "" {
        panic(fmt.Sprintf("required environment variable %s is not set", key))
    }
    return value
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
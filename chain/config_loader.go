package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Version  string         `json:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	overrideString := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	overrideInt := func(field *int, envVar string) error {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid integer value for %s: %s", envVar, val)
			}
			*field = intVal
		}
		return nil
	}

	overrideBool := func(field *bool, envVar string) error {
		if val := os.Getenv(envVar); val != "" {
			lowerVal := strings.ToLower(val)
			if lowerVal == "true" || lowerVal == "1" || lowerVal == "yes" {
				*field = true
			} else if lowerVal == "false" || lowerVal == "0" || lowerVal == "no" {
				*field = false
			} else {
				return fmt.Errorf("invalid boolean value for %s: %s", envVar, val)
			}
		}
		return nil
	}

	overrideString(&config.Database.Host, "DB_HOST")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")
	overrideString(&config.Database.SSLMode, "DB_SSL_MODE")
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")

	if err := overrideInt(&config.Database.Port, "DB_PORT"); err != nil {
		return err
	}
	if err := overrideInt(&config.Server.Port, "SERVER_PORT"); err != nil {
		return err
	}
	if err := overrideInt(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT"); err != nil {
		return err
	}
	if err := overrideInt(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT"); err != nil {
		return err
	}
	if err := overrideBool(&config.Server.DebugMode, "SERVER_DEBUG"); err != nil {
		return err
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if config.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return fmt.Errorf("server read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return fmt.Errorf("server write timeout cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
	}

	return nil
}
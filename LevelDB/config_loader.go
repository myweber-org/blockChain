package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	Environment  string `json:"environment"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	LogLevel string         `json:"log_level"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty")
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := json.Unmarshal(fileData, &config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	setDefaults(&config)
	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.Environment != "development" && config.Server.Environment != "production" {
		return errors.New("environment must be either 'development' or 'production'")
	}
	return nil
}

func setDefaults(config *AppConfig) {
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}
}
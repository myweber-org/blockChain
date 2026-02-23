
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	port, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port

	cfg.DBHost = getEnvString("DB_HOST", "localhost")

	dbPort, err := getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort

	debug, err := getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug

	apiKeysStr := getEnvString("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolValue, nil
	}
	return defaultValue, nil
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.DBHost == "" {
		return errors.New("database host cannot be empty")
	}

	return nil
}
package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
}

func LoadConfig() (*AppConfig, error) {
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	debugMode, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		ServerPort: port,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DebugMode:  debugMode,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Version  string         `json:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := json.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	envVars := map[string]string{
		"DB_HOST":       &config.Database.Host,
		"DB_PORT":       fmt.Sprintf("%d", config.Database.Port),
		"DB_USER":       &config.Database.Username,
		"DB_PASS":       &config.Database.Password,
		"DB_NAME":       &config.Database.Database,
		"SERVER_PORT":   fmt.Sprintf("%d", config.Server.Port),
		"READ_TIMEOUT":  fmt.Sprintf("%d", config.Server.ReadTimeout),
		"WRITE_TIMEOUT": fmt.Sprintf("%d", config.Server.WriteTimeout),
		"DEBUG_MODE":    fmt.Sprintf("%t", config.Server.DebugMode),
		"LOG_LEVEL":     &config.Server.LogLevel,
	}

	for envKey, defaultValue := range envVars {
		if envValue := os.Getenv(envKey); envValue != "" {
			switch envKey {
			case "DB_HOST", "DB_USER", "DB_PASS", "DB_NAME", "LOG_LEVEL":
				*envVars[envKey].(*string) = envValue
			case "DB_PORT", "SERVER_PORT", "READ_TIMEOUT", "WRITE_TIMEOUT":
				var intVal int
				if _, err := fmt.Sscanf(envValue, "%d", &intVal); err == nil {
					*envVars[envKey].(*int) = intVal
				}
			case "DEBUG_MODE":
				*envVars[envKey].(*bool) = (envValue == "true" || envValue == "1")
			}
		}
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	if config.Server.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}
	if config.Server.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
	}

	return nil
}

func (c *AppConfig) String() string {
	maskedConfig := *c
	maskedConfig.Database.Password = "******"
	data, _ := json.MarshalIndent(maskedConfig, "", "  ")
	return string(data)
}
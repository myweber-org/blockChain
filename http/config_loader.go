package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    AllowedIPs []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    var err error

    cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }

    cfg.DBHost = getStringEnv("DB_HOST", "localhost")
    
    cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    cfg.DebugMode, err = getBoolEnv("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }

    cfg.AllowedIPs = getStringSliceEnv("ALLOWED_IPS", []string{"127.0.0.1"})

    return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.Atoi(value)
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.ParseBool(value)
    }
    return defaultValue, nil
}

func getStringSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}package config

import (
	"encoding/json"
	"errors"
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
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Version  string         `json:"version"`
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config

	if configPath == "" {
		configPath = "config.json"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(fileData, &config); err != nil {
			return nil, err
		}
	}

	if err := loadEnvOverrides(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadEnvOverrides(config *Config) error {
	overrideStringField := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	overrideIntField := func(field *int, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
				*field = intVal
			}
		}
	}

	overrideBoolField := func(field *bool, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = strings.ToLower(val) == "true" || val == "1"
		}
	}

	overrideStringField(&config.Database.Host, "DB_HOST")
	overrideIntField(&config.Database.Port, "DB_PORT")
	overrideStringField(&config.Database.Username, "DB_USER")
	overrideStringField(&config.Database.Password, "DB_PASS")
	overrideStringField(&config.Database.Database, "DB_NAME")

	overrideIntField(&config.Server.Port, "SERVER_PORT")
	overrideIntField(&config.Server.ReadTimeout, "READ_TIMEOUT")
	overrideIntField(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
	overrideBoolField(&config.Server.DebugMode, "DEBUG_MODE")
	overrideStringField(&config.Server.LogLevel, "LOG_LEVEL")

	return nil
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Database.Username == "" {
		return errors.New("database username is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
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
		"fatal": true,
	}
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return errors.New("invalid log level")
	}

	return nil
}
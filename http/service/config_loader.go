package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Port        int
    DatabaseURL string
    LogLevel    string
    MaxWorkers  int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Port:        getEnvAsInt("APP_PORT", 8080),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        MaxWorkers:  getEnvAsInt("MAX_WORKERS", 10),
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
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return value
}

func validateConfig(cfg *Config) error {
    if cfg.Port < 1 || cfg.Port > 65535 {
        return &ConfigError{Field: "Port", Message: "port must be between 1 and 65535"}
    }
    if !strings.HasPrefix(cfg.DatabaseURL, "postgres://") {
        return &ConfigError{Field: "DatabaseURL", Message: "database URL must start with postgres://"}
    }
    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[cfg.LogLevel] {
        return &ConfigError{Field: "LogLevel", Message: "invalid log level"}
    }
    if cfg.MaxWorkers < 1 || cfg.MaxWorkers > 100 {
        return &ConfigError{Field: "MaxWorkers", Message: "max workers must be between 1 and 100"}
    }
    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        configPath = "config.yaml"
    }

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Database.Host == "" {
        config.Database.Host = "localhost"
    }

    if config.Database.Port == 0 {
        config.Database.Port = 5432
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    if config.Server.Mode == "" {
        config.Server.Mode = "development"
    }

    if config.LogLevel == "" {
        config.LogLevel = "info"
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }

    if config.Database.Username == "" {
        return fmt.Errorf("database username is required")
    }

    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    validModes := map[string]bool{
        "development": true,
        "staging":     true,
        "production":  true,
    }

    if !validModes[config.Server.Mode] {
        return fmt.Errorf("invalid server mode: %s", config.Server.Mode)
    }

    return nil
}
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var err error

	cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode, err = getBoolEnv("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	cfg.DatabaseURL = getStringEnv("DATABASE_URL", "postgres://localhost:5432/appdb")
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL cannot be empty")
	}

	cfg.CacheTTL, err = getIntEnv("CACHE_TTL", 300)
	if err != nil {
		return nil, err
	}

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
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
	if value := os.Getenv(key); value != "" {
		lowerValue := strings.ToLower(value)
		if lowerValue == "true" || lowerValue == "1" {
			return true, nil
		}
		if lowerValue == "false" || lowerValue == "0" {
			return false, nil
		}
		return false, errors.New("invalid boolean value for " + key)
	}
	return defaultValue, nil
}package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	cfg.Server.Host = getEnvOrDefault("SERVER_HOST", cfg.Server.Host)
	cfg.Server.Port = getEnvIntOrDefault("SERVER_PORT", cfg.Server.Port)

	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvIntOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.Name = getEnvOrDefault("DB_NAME", cfg.Database.Name)
	cfg.Database.User = getEnvOrDefault("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)

	cfg.Logging.Level = getEnvOrDefault("LOG_LEVEL", cfg.Logging.Level)
	cfg.Logging.Output = getEnvOrDefault("LOG_OUTPUT", cfg.Logging.Output)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
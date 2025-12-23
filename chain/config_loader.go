package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    DebugMode    bool   `yaml:"debug" env:"DEBUG_MODE"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)
    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    overrideStruct(config)
}

func overrideStruct(v interface{}) {
    // Implementation would use reflection to check struct tags
    // and override values from environment variables
    // Simplified for this example
}package config

import (
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
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
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}
	
	overrideWithEnvVars(config)
	
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	
	return config, nil
}

func overrideWithEnvVars(config *Config) {
	overrideStruct(&config.Database)
	overrideStruct(&config.Server)
}

func overrideStruct(target interface{}) {
	// Implementation would use reflection to read struct tags
	// and override values from environment variables
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return NewValidationError("database host is required")
	}
	
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return NewValidationError("database port must be between 1 and 65535")
	}
	
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return NewValidationError("server port must be between 1 and 65535")
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return NewValidationError("invalid log level")
	}
	
	return nil
}

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func (e *ValidationError) Error() string {
	return "config validation error: " + e.Message
}package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port

	cfg.DBHost = getEnv("DB_HOST", "localhost")

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort

	debug, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug

	keys := strings.Split(getEnv("API_KEYS", ""), ",")
	cfg.APIKeys = keys

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}package config

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
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    // Database configuration
    cfg.Database.Host = getEnv("DB_HOST", "localhost")
    cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
    cfg.Database.Username = getEnv("DB_USER", "postgres")
    cfg.Database.Password = getEnv("DB_PASSWORD", "")
    cfg.Database.Database = getEnv("DB_NAME", "appdb")
    cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

    // Server configuration
    cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
    cfg.Server.ReadTimeout = getEnvAsInt("READ_TIMEOUT", 30)
    cfg.Server.WriteTimeout = getEnvAsInt("WRITE_TIMEOUT", 30)
    cfg.Server.DebugMode = getEnvAsBool("DEBUG_MODE", false)

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
        return &ConfigValidationError{Errors: errors}
    }
    return nil
}

type ConfigValidationError struct {
    Errors []string
}

func (e *ConfigValidationError) Error() string {
    return "configuration validation failed: " + strings.Join(e.Errors, ", ")
}
package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("DATABASE_URL must start with postgres://")
    }
    cfg.DatabaseURL = dbURL
    
    cacheStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cacheEnabled, err := strconv.ParseBool(cacheStr)
    if err != nil {
        return nil, fmt.Errorf("invalid CACHE_ENABLED value: %v", err)
    }
    cfg.CacheEnabled = cacheEnabled
    
    logLevel := getEnvWithDefault("LOG_LEVEL", "info")
    validLevels := map[string]bool{
        "debug": true,
        "info": true,
        "warn": true,
        "error": true,
    }
    if !validLevels[logLevel] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", logLevel)
    }
    cfg.LogLevel = logLevel
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port" env:"SERVER_PORT" default:"8080"`
	DBHost     string `json:"db_host" env:"DB_HOST" default:"localhost"`
	DBPort     int    `json:"db_port" env:"DB_PORT" default:"5432"`
	DBName     string `json:"db_name" env:"DB_NAME" default:"appdb"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE" default:"false"`
	LogLevel   string `json:"log_level" env:"LOG_LEVEL" default:"info"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}
	
	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}
	
	if err := applyDefaults(cfg); err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}
	
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cfg, nil
}

func loadFromFile(path string, cfg *Config) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, cfg)
}

func loadFromEnv(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		envTag := structField.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		
		if err := setFieldValue(field, envValue); err != nil {
			return fmt.Errorf("invalid value for %s: %w", envTag, err)
		}
	}
	
	return nil
}

func applyDefaults(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		if !field.IsZero() {
			continue
		}
		
		defaultTag := structField.Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		
		if err := setFieldValue(field, defaultTag); err != nil {
			return fmt.Errorf("invalid default for %s: %w", structField.Name, err)
		}
	}
	
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func validateConfig(cfg *Config) error {
	var validationErrors []string
	
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		validationErrors = append(validationErrors, "server_port must be between 1 and 65535")
	}
	
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		validationErrors = append(validationErrors, "db_port must be between 1 and 65535")
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		validationErrors = append(validationErrors, "log_level must be one of: debug, info, warn, error")
	}
	
	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
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

	apiKeys := getEnvString("API_KEYS", "")
	if apiKeys != "" {
		cfg.APIKeys = strings.Split(apiKeys, ",")
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intVal, nil
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value := os.Getenv(key); value != "" {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolVal, nil
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
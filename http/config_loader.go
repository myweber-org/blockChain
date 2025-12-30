package config

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DebugMode:  getEnvAsBool("DEBUG_MODE", false),
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
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DBName     string `json:"db_name" env:"DB_NAME"`
	LogLevel   string `json:"log_level" env:"LOG_LEVEL"`
	CacheTTL   int    `json:"cache_ttl" env:"CACHE_TTL"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config
	
	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("invalid config path: %w", err)
		}
		
		file, err := os.Open(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()
		
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&cfg); err != nil {
			return nil, fmt.Errorf("failed to decode config: %w", err)
		}
	}
	
	if err := loadEnvVars(&cfg); err != nil {
		return nil, err
	}
	
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

func loadEnvVars(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		
		if envTag == "" {
			continue
		}
		
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		
		fieldValue := v.Field(i)
		
		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Int:
			var intVal int
			if _, err := fmt.Sscanf(envValue, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid integer value for %s: %s", envTag, envValue)
			}
			fieldValue.SetInt(int64(intVal))
		}
	}
	
	return nil
}

func validateConfig(cfg *Config) error {
	var errors []string
	
	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		errors = append(errors, "server_port must be between 1 and 65535")
	}
	
	if cfg.DBHost == "" {
		errors = append(errors, "db_host is required")
	}
	
	if cfg.DBPort <= 0 || cfg.DBPort > 65535 {
		errors = append(errors, "db_port must be between 1 and 65535")
	}
	
	if cfg.DBName == "" {
		errors = append(errors, "db_name is required")
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		errors = append(errors, "log_level must be one of: debug, info, warn, error")
	}
	
	if cfg.CacheTTL < 0 {
		errors = append(errors, "cache_ttl cannot be negative")
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(errors, "; "))
	}
	
	return nil
}package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	DebugMode    bool   `yaml:"debug_mode"`
	LogLevel     string `yaml:"log_level"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	if filePath == "" {
		return nil, errors.New("config file path cannot be empty")
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
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
	}
	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level, must be debug, info, warn, or error")
	}

	if config.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if config.Database.Database == "" {
		return errors.New("database name cannot be empty")
	}

	return nil
}
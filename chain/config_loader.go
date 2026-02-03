package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&cfg.Database)
    overrideFromEnv(&cfg.Server)

    return &cfg, nil
}

func overrideFromEnv(config interface{}) {
    // Implementation would use reflection to check struct tags
    // and override values from environment variables
    // Simplified placeholder implementation
}

func DefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.yaml"
    }
    return filepath.Join(homeDir, ".app", "config.yaml")
}
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		defaultVal := field.Tag.Get("default")

		envValue := os.Getenv(envTag)
		if envValue == "" {
			envValue = defaultVal
		}

		if envValue == "" {
			return nil, fmt.Errorf("environment variable %s is required", envTag)
		}

		fieldValue := v.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Int:
			intVal, err := strconv.Atoi(envValue)
			if err != nil {
				return nil, fmt.Errorf("invalid integer value for %s: %v", envTag, err)
			}
			fieldValue.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(strings.ToLower(envValue))
			if err != nil {
				return nil, fmt.Errorf("invalid boolean value for %s: %v", envTag, err)
			}
			fieldValue.SetBool(boolVal)
		default:
			return nil, fmt.Errorf("unsupported field type: %s", field.Type.Kind())
		}
	}

	return cfg, nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required,min=1,max=65535"`
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  int    `yaml:"read_timeout" validate:"required,min=1"`
	WriteTimeout int    `yaml:"write_timeout" validate:"required,min=1"`
	Environment  string `yaml:"environment" validate:"required,oneof=development staging production"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Features struct {
		EnableCache bool `yaml:"enable_cache"`
		MaxWorkers  int  `yaml:"max_workers" validate:"required,min=1"`
	} `yaml:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty")
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if config.Features.MaxWorkers < 1 {
		return errors.New("max workers must be at least 1")
	}

	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}

	if !validEnvs[config.Server.Environment] {
		return errors.New("invalid environment specified")
	}

	return nil
}

func SaveConfig(config *AppConfig, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    Debug        bool   `yaml:"debug" env:"DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
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
    // This function would iterate through struct fields
    // and override values from environment variables
    // For brevity, implementing simple overrides
    if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
        var port int
        if _, err := fmt.Sscanf(envPort, "%d", &port); err == nil {
            config.Server.Port = port
        }
    }

    if envDebug := os.Getenv("DEBUG"); envDebug != "" {
        config.Server.Debug = envDebug == "true" || envDebug == "1"
    }

    if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
        config.LogLevel = envLogLevel
    }
}

func DefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.yaml"
    }
    return filepath.Join(homeDir, ".app", "config.yaml")
}
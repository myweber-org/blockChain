package config

import (
    "errors"
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.Server.Host == "" {
        return errors.New("server host cannot be empty")
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return errors.New("server port must be between 1 and 65535")
    }
    if cfg.Database.Host == "" {
        return errors.New("database host cannot be empty")
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    if !validLogLevels[cfg.LogLevel] {
        return errors.New("invalid log level")
    }

    return nil
}package config

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
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Features struct {
		CacheEnabled  bool `yaml:"cache_enabled" env:"CACHE_ENABLED"`
		Metrics       bool `yaml:"metrics" env:"METRICS_ENABLED"`
		RateLimit     int  `yaml:"rate_limit" env:"RATE_LIMIT"`
		MaxFileSizeMB int  `yaml:"max_file_size_mb" env:"MAX_FILE_SIZE_MB"`
	} `yaml:"features"`
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

	if err := applyEnvOverrides(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func applyEnvOverrides(config *AppConfig) error {
	overrideString := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	overrideInt := func(field *int, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
				*field = intVal
			}
		}
	}

	overrideBool := func(field *bool, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val == "true" || val == "1" || val == "yes"
		}
	}

	overrideString(&config.Database.Host, "DB_HOST")
	overrideInt(&config.Database.Port, "DB_PORT")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")

	overrideInt(&config.Server.Port, "SERVER_PORT")
	overrideInt(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
	overrideInt(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
	overrideBool(&config.Server.DebugMode, "SERVER_DEBUG")
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")

	overrideBool(&config.Features.CacheEnabled, "CACHE_ENABLED")
	overrideBool(&config.Features.Metrics, "METRICS_ENABLED")
	overrideInt(&config.Features.RateLimit, "RATE_LIMIT")
	overrideInt(&config.Features.MaxFileSizeMB, "MAX_FILE_SIZE_MB")

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if config.Database.Host == "" {
		return errors.New("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if config.Database.Database == "" {
		return errors.New("database name is required")
	}

	if config.Features.RateLimit < 0 {
		return errors.New("rate limit cannot be negative")
	}

	if config.Features.MaxFileSizeMB < 0 {
		return errors.New("max file size cannot be negative")
	}

	return nil
}

func SaveDefaultConfig(path string) error {
	defaultConfig := AppConfig{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
			DebugMode:    false,
			LogLevel:     "info",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "",
			Database: "appdb",
		},
	}

	defaultConfig.Features.CacheEnabled = true
	defaultConfig.Features.Metrics = false
	defaultConfig.Features.RateLimit = 100
	defaultConfig.Features.MaxFileSizeMB = 10

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
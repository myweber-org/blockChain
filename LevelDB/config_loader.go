package config

import (
	"errors"
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
		URL      string `yaml:"url" env:"DB_URL"`
		MaxConns int    `yaml:"max_connections" env:"DB_MAX_CONNS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := overrideFromEnv(&cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) error {
	if val := os.Getenv("SERVER_HOST"); val != "" {
		cfg.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			return errors.New("invalid SERVER_PORT value")
		}
		cfg.Server.Port = port
	}
	if val := os.Getenv("DB_URL"); val != "" {
		cfg.Database.URL = val
	}
	if val := os.Getenv("DB_MAX_CONNS"); val != "" {
		maxConns, err := strconv.Atoi(val)
		if err != nil {
			return errors.New("invalid DB_MAX_CONNS value")
		}
		cfg.Database.MaxConns = maxConns
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		cfg.LogLevel = val
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if cfg.Database.URL == "" {
		return errors.New("database URL is required")
	}
	if cfg.Database.MaxConns < 1 {
		return errors.New("database max connections must be at least 1")
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
}
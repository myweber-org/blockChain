package config

import (
	"os"
	"path/filepath"
	"strings"

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
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Format string `yaml:"format" env:"LOG_FORMAT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

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

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(cfg)

	return cfg, nil
}

func overrideFromEnv(cfg *Config) {
	overrideField(&cfg.Server.Host, "SERVER_HOST")
	overrideField(&cfg.Server.Port, "SERVER_PORT")
	overrideField(&cfg.Database.Host, "DB_HOST")
	overrideField(&cfg.Database.Port, "DB_PORT")
	overrideField(&cfg.Database.Name, "DB_NAME")
	overrideField(&cfg.Database.User, "DB_USER")
	overrideField(&cfg.Database.Password, "DB_PASSWORD")
	overrideField(&cfg.Database.SSLMode, "DB_SSL_MODE")
	overrideField(&cfg.Logging.Level, "LOG_LEVEL")
	overrideField(&cfg.Logging.Format, "LOG_FORMAT")
}

func overrideField(field interface{}, envVar string) {
	envValue := os.Getenv(envVar)
	if envValue == "" {
		return
	}

	switch v := field.(type) {
	case *string:
		*v = envValue
	case *int:
		if intVal, err := parseInt(envValue); err == nil {
			*v = intVal
		}
	}
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
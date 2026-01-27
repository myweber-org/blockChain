package config

import (
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
        Driver   string `yaml:"driver"`
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}package config

import (
	"errors"
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

	port, err := getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port

	cfg.DBHost = getStringEnv("DB_HOST", "localhost")

	dbPort, err := getIntEnv("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort

	debug, err := getBoolEnv("DEBUG_MODE")
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug

	apiKeys := getStringEnv("API_KEYS", "")
	if apiKeys != "" {
		cfg.APIKeys = strings.Split(apiKeys, ",")
	}

	if err := validateConfig(cfg); err != nil {
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

func getBoolEnv(key string) (bool, error) {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolValue, nil
	}
	return false, nil
}

func validateConfig(cfg *Config) error {
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
package config

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
	FeatureFlags map[string]bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		FeatureFlags: make(map[string]bool),
	}

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

	flags := strings.Split(getEnv("FEATURE_FLAGS", ""), ",")
	for _, flag := range flags {
		trimmed := strings.TrimSpace(flag)
		if trimmed != "" {
			cfg.FeatureFlags[trimmed] = true
		}
	}

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
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    ServerPort string `json:"server_port"`
    DBHost     string `json:"db_host"`
    DBPort     string `json:"db_port"`
    DebugMode  bool   `json:"debug_mode"`
}

func LoadConfig(configPath string) (*Config, error) {
    if configPath == "" {
        configPath = getDefaultConfigPath()
    }

    fileData, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := json.Unmarshal(fileData, &config); err != nil {
        return nil, err
    }

    overrideFromEnv(&config)
    return &config, nil
}

func getDefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.json"
    }
    return filepath.Join(homeDir, ".app", "config.json")
}

func overrideFromEnv(cfg *Config) {
    if port := os.Getenv("SERVER_PORT"); port != "" {
        cfg.ServerPort = port
    }
    if host := os.Getenv("DB_HOST"); host != "" {
        cfg.DBHost = host
    }
    if port := os.Getenv("DB_PORT"); port != "" {
        cfg.DBPort = port
    }
    if debug := os.Getenv("DEBUG_MODE"); debug == "true" {
        cfg.DebugMode = true
    }
}package config

import (
	"errors"
	"io"
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

func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
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

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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

	cfg.loadFromEnv()
	return &cfg, nil
}

func (c *Config) loadFromEnv() {
	c.Server.Host = getEnvOrDefault("SERVER_HOST", c.Server.Host)
	c.Server.Port = getEnvIntOrDefault("SERVER_PORT", c.Server.Port)
	c.Database.Host = getEnvOrDefault("DB_HOST", c.Database.Host)
	c.Database.Port = getEnvIntOrDefault("DB_PORT", c.Database.Port)
	c.Database.Name = getEnvOrDefault("DB_NAME", c.Database.Name)
	c.Database.User = getEnvOrDefault("DB_USER", c.Database.User)
	c.Database.Password = getEnvOrDefault("DB_PASSWORD", c.Database.Password)
	c.Logging.Level = getEnvOrDefault("LOG_LEVEL", c.Logging.Level)
	c.Logging.Output = getEnvOrDefault("LOG_OUTPUT", c.Logging.Output)
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
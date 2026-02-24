package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort   int
    DatabaseURL  string
    LogLevel     string
    CacheEnabled bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:     getEnv("LOG_LEVEL", "info"),
        CacheEnabled: getEnvAsBool("CACHE_ENABLED", true),
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
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }

    if !strings.HasPrefix(cfg.DatabaseURL, "postgres://") {
        return &ConfigError{Field: "DatabaseURL", Message: "must be a PostgreSQL connection string"}
    }

    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return &ConfigError{Field: "LogLevel", Message: "must be one of: debug, info, warn, error"}
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
	"io/ioutil"
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
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	overrideFromEnv(&cfg)

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	overrideString(&cfg.Server.Host, "SERVER_HOST")
	overrideInt(&cfg.Server.Port, "SERVER_PORT")
	overrideString(&cfg.Database.Host, "DB_HOST")
	overrideInt(&cfg.Database.Port, "DB_PORT")
	overrideString(&cfg.Database.Name, "DB_NAME")
	overrideString(&cfg.Database.Username, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASS")
	overrideString(&cfg.Logging.Level, "LOG_LEVEL")
	overrideString(&cfg.Logging.Output, "LOG_OUTPUT")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}package config

import (
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort   int
    DatabaseURL  string
    CacheEnabled bool
    LogLevel     string
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:   8080,
        DatabaseURL:  "localhost:5432",
        CacheEnabled: true,
        LogLevel:     "info",
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
            cfg.ServerPort = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheFlag := os.Getenv("ENABLE_CACHE"); cacheFlag != "" {
        cfg.CacheEnabled = strings.ToLower(cacheFlag) == "true"
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        validLevels := map[string]bool{
            "debug": true,
            "info":  true,
            "warn":  true,
            "error": true,
        }
        if validLevels[strings.ToLower(logLevel)] {
            cfg.LogLevel = strings.ToLower(logLevel)
        }
    }

    return cfg, nil
}package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := make(Config)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format on line %d", lineNumber)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = os.ExpandEnv(value)
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

func (c Config) Get(key string) (string, error) {
	value, exists := c[key]
	if !exists {
		return "", fmt.Errorf("key '%s' not found in configuration", key)
	}
	return value, nil
}
package config

import (
	"io/ioutil"
	"os"

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
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
		Name     string `yaml:"name" env:"DB_NAME"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}

	if filePath != "" {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	overrideFromEnv(config)
	return config, nil
}

func overrideFromEnv(c *Config) {
	if val := os.Getenv("SERVER_HOST"); val != "" {
		c.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		port := 0
		fmt.Sscanf(val, "%d", &port)
		if port > 0 {
			c.Server.Port = port
		}
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		c.Database.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		port := 0
		fmt.Sscanf(val, "%d", &port)
		if port > 0 {
			c.Database.Port = port
		}
	}
	if val := os.Getenv("DB_USER"); val != "" {
		c.Database.Username = val
	}
	if val := os.Getenv("DB_PASS"); val != "" {
		c.Database.Password = val
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		c.Database.Name = val
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.LogLevel = val
	}
}package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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
	Environment  string `yaml:"environment"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	if filePath == "" {
		return nil, errors.New("config file path cannot be empty")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
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

	return nil
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvOrDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    debugStr := getEnvOrDefault("ENABLE_DEBUG", "false")
    cfg.EnableDebug = strings.ToLower(debugStr) == "true"
    
    originsStr := getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return ErrInvalidPort
    }
    
    if cfg.DatabaseURL == "" {
        return ErrMissingDatabaseURL
    }
    
    return nil
}

var (
    ErrInvalidPort = newConfigError("invalid server port")
    ErrMissingDatabaseURL = newConfigError("database URL is required")
)

type ConfigError struct {
    message string
}

func newConfigError(msg string) *ConfigError {
    return &ConfigError{message: msg}
}

func (e *ConfigError) Error() string {
    return "config error: " + e.message
}
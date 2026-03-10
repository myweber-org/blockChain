package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
    AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        CacheEnabled:   true,
        MaxConnections: 100,
        AllowedOrigins: []string{"http://localhost:3000"},
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
            cfg.ServerPort = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheStr := os.Getenv("CACHE_ENABLED"); cacheStr != "" {
        cfg.CacheEnabled = strings.ToLower(cacheStr) == "true"
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        if maxConn, err := strconv.Atoi(maxConnStr); err == nil && maxConn > 0 {
            cfg.MaxConnections = maxConn
        }
    }

    if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
        cfg.AllowedOrigins = strings.Split(origins, ",")
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return ErrInvalidPort
    }

    if cfg.DatabaseURL == "" {
        return ErrMissingDatabaseURL
    }

    if cfg.MaxConnections < 1 {
        return ErrInvalidMaxConnections
    }

    return nil
}

var (
    ErrInvalidPort           = errors.New("invalid server port")
    ErrMissingDatabaseURL    = errors.New("database URL is required")
    ErrInvalidMaxConnections = errors.New("max connections must be positive")
)
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
	AllowedIPs []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

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

	ips := getEnv("ALLOWED_IPS", "127.0.0.1,::1")
	cfg.AllowedIPs = strings.Split(ips, ",")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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
	DebugMode  bool
	APIKey     string
	Timeout    int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port

	debugStr := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"

	cfg.APIKey = os.Getenv("API_KEY")
	if cfg.APIKey == "" {
		return nil, errors.New("API_KEY is required")
	}

	timeoutStr := os.Getenv("TIMEOUT_SECONDS")
	if timeoutStr == "" {
		timeoutStr = "30"
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout < 1 {
		return nil, errors.New("invalid TIMEOUT_SECONDS value")
	}
	cfg.Timeout = timeout

	return cfg, nil
}package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	ServerPort int    `json:"server_port"`
	LogLevel   string `json:"log_level"`
	CacheSize  int    `json:"cache_size"`
	EnableTLS  bool   `json:"enable_tls"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = filepath.Join("config", "settings.json")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	setDefaults(&config)
	return &config, nil
}

func validateConfig(c *Config) error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return errors.New("invalid log level")
	}

	if c.CacheSize < 0 {
		return errors.New("cache size cannot be negative")
	}

	return nil
}

func setDefaults(c *Config) {
	if c.ServerPort == 0 {
		c.ServerPort = 8080
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.CacheSize == 0 {
		c.CacheSize = 100
	}
}
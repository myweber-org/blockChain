package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	portStr := getEnv("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port
	
	debugStr := getEnv("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
	
	hostsStr := getEnv("ALLOWED_HOSTS", "localhost,127.0.0.1")
	cfg.AllowedHosts = strings.Split(hostsStr, ",")
	
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
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort    int
	DatabaseURL   string
	LogLevel      string
	CacheEnabled  bool
	MaxConnections int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	var err error
	cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
	cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	cfg.CacheEnabled = getBoolEnv("CACHE_ENABLED", true)
	cfg.MaxConnections, err = getIntEnv("MAX_CONNECTIONS", 100)
	if err != nil {
		return nil, err
	}
	
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return nil, errors.New("invalid server port range")
	}
	
	if cfg.MaxConnections < 1 {
		return nil, errors.New("max connections must be positive")
	}
	
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue, nil
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, errors.New("invalid integer value for " + key)
	}
	return value, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	valueStr := strings.ToLower(getEnv(key, ""))
	if valueStr == "" {
		return defaultValue
	}
	
	return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}
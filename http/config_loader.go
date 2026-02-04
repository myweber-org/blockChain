
package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    cacheEnabledStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cacheEnabled, err := strconv.ParseBool(cacheEnabledStr)
    if err != nil {
        return nil, fmt.Errorf("invalid CACHE_ENABLED value: %v", err)
    }
    cfg.CacheEnabled = cacheEnabled
    
    cfg.LogLevel = strings.ToUpper(getEnvWithDefault("LOG_LEVEL", "INFO"))
    validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
    if !validLogLevels[cfg.LogLevel] {
        return nil, fmt.Errorf("invalid LOG_LEVEL value: %s", cfg.LogLevel)
    }
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
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
	"sync"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DebugMode  bool   `json:"debug_mode"`
	CacheTTL   int    `json:"cache_ttl"`
}

var (
	config     *AppConfig
	configOnce sync.Once
)

func LoadConfig() (*AppConfig, error) {
	var err error
	configOnce.Do(func() {
		configPath := getConfigPath()
		config, err = loadFromFile(configPath)
		if err != nil {
			config = loadFromEnv()
		}
	})
	return config, err
}

func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return filepath.Join(".", "config.json")
}

func loadFromFile(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadFromEnv() *AppConfig {
	return &AppConfig{
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
		DebugMode:  getEnvBoolOrDefault("DEBUG_MODE", false),
		CacheTTL:   getEnvIntOrDefault("CACHE_TTL", 300),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
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
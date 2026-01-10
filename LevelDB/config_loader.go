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
}
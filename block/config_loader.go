package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if port := os.Getenv("APP_PORT"); port != "" {
		config.ServerPort = port
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.DBHost = dbHost
	}

	return &config, nil
}

func SaveConfig(config *AppConfig, configPath string) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	LogLevel   string `json:"log_level"`
}

var (
	config     *Config
	configOnce sync.Once
)

func Load() *Config {
	configOnce.Do(func() {
		config = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			LogLevel:   getEnv("LOG_LEVEL", "info"),
		}

		if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
			loadFromFile(configFile)
		}
	})
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadFromFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return
	}

	if fileConfig.ServerPort != "" {
		config.ServerPort = fileConfig.ServerPort
	}
	if fileConfig.DBHost != "" {
		config.DBHost = fileConfig.DBHost
	}
	if fileConfig.DBPort != "" {
		config.DBPort = fileConfig.DBPort
	}
	if fileConfig.LogLevel != "" {
		config.LogLevel = fileConfig.LogLevel
	}
}
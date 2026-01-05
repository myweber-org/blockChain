package config

import (
	"encoding/json"
	"os"
	"sync"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

var (
	config     *AppConfig
	configOnce sync.Once
)

func LoadConfig() *AppConfig {
	configOnce.Do(func() {
		configFile := os.Getenv("CONFIG_FILE")
		if configFile == "" {
			configFile = "config.json"
		}

		data, err := os.ReadFile(configFile)
		if err != nil {
			config = &AppConfig{
				ServerPort: getEnv("SERVER_PORT", "8080"),
				DBHost:     getEnv("DB_HOST", "localhost"),
				DBPort:     5432,
				DebugMode:  getEnv("DEBUG_MODE", "false") == "true",
			}
			return
		}

		var cfg AppConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			panic("Failed to parse config file: " + err.Error())
		}

		config = &cfg
	})

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
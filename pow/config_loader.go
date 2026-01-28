package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	APIKey      string
	LogLevel    string
	Port        string
}

func LoadConfig(filePath string) (*Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	config := &Config{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = substituteEnvVars(value)

		switch key {
		case "DATABASE_URL":
			config.DatabaseURL = value
		case "API_KEY":
			config.APIKey = value
		case "LOG_LEVEL":
			config.LogLevel = value
		case "PORT":
			config.Port = value
		}
	}

	return config, nil
}

func substituteEnvVars(value string) string {
	return os.ExpandEnv(value)
}
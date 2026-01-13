package config

import (
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  bool   `json:"ssl_mode" env:"DB_SSL"`
}

type AppConfig struct {
	Debug      bool           `json:"debug" env:"APP_DEBUG"`
	Port       int            `json:"port" env:"APP_PORT"`
	SecretKey  string         `json:"secret_key" env:"APP_SECRET"`
	Database   DatabaseConfig `json:"database"`
	AllowedIPs []string       `json:"allowed_ips" env:"APP_ALLOWED_IPS"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{
		Port: 8080,
		Database: DatabaseConfig{
			Port: 5432,
		},
	}

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	loadFromEnv(config)
	return config, validateConfig(config)
}

func loadFromEnv(config *AppConfig) {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	config.Debug = getEnvBool("APP_DEBUG", config.Debug, envMap)
	config.Port = getEnvInt("APP_PORT", config.Port, envMap)
	config.SecretKey = getEnvString("APP_SECRET", config.SecretKey, envMap)
	
	config.Database.Host = getEnvString("DB_HOST", config.Database.Host, envMap)
	config.Database.Port = getEnvInt("DB_PORT", config.Database.Port, envMap)
	config.Database.Username = getEnvString("DB_USER", config.Database.Username, envMap)
	config.Database.Password = getEnvString("DB_PASS", config.Database.Password, envMap)
	config.Database.SSLMode = getEnvBool("DB_SSL", config.Database.SSLMode, envMap)
	
	if ips := getEnvString("APP_ALLOWED_IPS", "", envMap); ips != "" {
		config.AllowedIPs = strings.Split(ips, ",")
	}
}

func validateConfig(config *AppConfig) error {
	if config.Port < 1 || config.Port > 65535 {
		return &ConfigError{Field: "port", Message: "port must be between 1 and 65535"}
	}
	
	if config.Database.Host == "" {
		return &ConfigError{Field: "database.host", Message: "database host is required"}
	}
	
	if config.SecretKey == "" {
		return &ConfigError{Field: "secret_key", Message: "secret key is required"}
	}
	
	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}

func getEnvString(key, defaultValue string, envMap map[string]string) string {
	if val, exists := envMap[key]; exists && val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int, envMap map[string]string) int {
	if val, exists := envMap[key]; exists && val != "" {
		var result int
		if _, err := fmt.Sscanf(val, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool, envMap map[string]string) bool {
	if val, exists := envMap[key]; exists && val != "" {
		switch strings.ToLower(val) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return defaultValue
}package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
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

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	cfg.Server.Host = getEnvOrDefault("SERVER_HOST", cfg.Server.Host)
	cfg.Server.Port = getEnvIntOrDefault("SERVER_PORT", cfg.Server.Port)

	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvIntOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.Name = getEnvOrDefault("DB_NAME", cfg.Database.Name)
	cfg.Database.User = getEnvOrDefault("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)
	cfg.Database.SSLMode = getEnvOrDefault("DB_SSL_MODE", cfg.Database.SSLMode)

	cfg.Logging.Level = getEnvOrDefault("LOG_LEVEL", cfg.Logging.Level)
	cfg.Logging.Output = getEnvOrDefault("LOG_OUTPUT", cfg.Logging.Output)
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
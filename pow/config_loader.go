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
	APIKeys    []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
		APIKeys:    getEnvAsSlice("API_KEYS", []string{}, ","),
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

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return ErrInvalidPort
	}
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return ErrInvalidPort
	}
	if len(cfg.APIKeys) == 0 {
		return ErrMissingAPIKeys
	}
	return nil
}

var (
	ErrInvalidPort    = ConfigError{"invalid port number"}
	ErrMissingAPIKeys = ConfigError{"at least one API key is required"}
)

type ConfigError struct {
	Message string
}

func (e ConfigError) Error() string {
	return "config error: " + e.Message
}package config

import (
    "errors"
    "io/ioutil"
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
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Debug    bool           `yaml:"debug"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.Database.Host == "" {
        return errors.New("database host cannot be empty")
    }
    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        return errors.New("database port must be between 1 and 65535")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return errors.New("server port must be between 1 and 65535")
    }
    return nil
}
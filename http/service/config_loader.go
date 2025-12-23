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
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}package config

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

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		DebugMode:   getEnvAsBool("DEBUG_MODE", false),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
		AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}, ","),
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
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
    config.Database.Port = getEnvIntOrDefault("DB_PORT", config.Database.Port)
    config.Database.Username = getEnvOrDefault("DB_USER", config.Database.Username)
    config.Database.Password = getEnvOrDefault("DB_PASS", config.Database.Password)
    config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)

    config.Server.Port = getEnvIntOrDefault("SERVER_PORT", config.Server.Port)
    config.Server.ReadTimeout = getEnvIntOrDefault("READ_TIMEOUT", config.Server.ReadTimeout)
    config.Server.WriteTimeout = getEnvIntOrDefault("WRITE_TIMEOUT", config.Server.WriteTimeout)
    config.Server.DebugMode = getEnvBoolOrDefault("DEBUG_MODE", config.Server.DebugMode)

    config.LogLevel = getEnvOrDefault("LOG_LEVEL", config.LogLevel)
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

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        return value == "true" || value == "1" || value == "yes"
    }
    return defaultValue
}

func DefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.yaml"
    }
    return filepath.Join(homeDir, ".app", "config.yaml")
}package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	DebugMode    bool
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	LogLevel string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}

	cfg.Database = DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", ""),
		Database: getEnv("DB_NAME", "appdb"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}

	readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT value")
	}

	writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT value")
	}

	debugMode, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		return nil, errors.New("invalid DEBUG_MODE value")
	}

	cfg.Server = ServerConfig{
		Port:         serverPort,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		DebugMode:    debugMode,
	}

	cfg.LogLevel = strings.ToUpper(getEnv("LOG_LEVEL", "INFO"))

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

func validateConfig(cfg *Config) error {
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	validLogLevels := map[string]bool{
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
		"FATAL": true,
	}

	if !validLogLevels[cfg.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func (c *Config) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid port number: %d", c.Server.Port)
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Database.Port)
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.Port)
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.ReadTimeout)
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.WriteTimeout)
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        config.Server.DebugMode = val == "true" || val == "1"
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }
}package config

import (
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort   int
    DatabaseURL  string
    LogLevel     string
    FeatureFlags map[string]bool
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:     getEnv("LOG_LEVEL", "info"),
        FeatureFlags: parseFeatureFlags(getEnv("FEATURE_FLAGS", "")),
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
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    
    intValue, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return intValue
}

func parseFeatureFlags(flagsStr string) map[string]bool {
    flags := make(map[string]bool)
    if flagsStr == "" {
        return flags
    }

    pairs := strings.Split(flagsStr, ",")
    for _, pair := range pairs {
        parts := strings.Split(pair, "=")
        if len(parts) == 2 {
            flags[strings.TrimSpace(parts[0])] = strings.ToLower(strings.TrimSpace(parts[1])) == "true"
        }
    }
    return flags
}

func validateConfig(cfg *AppConfig) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }
    
    if cfg.DatabaseURL == "" {
        return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
    }
    
    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return &ConfigError{Field: "LogLevel", Message: "invalid log level"}
    }
    
    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}
package config

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
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Debug    bool
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}

	config.Database = DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", ""),
		Database: getEnv("DB_NAME", "appdb"),
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}

	config.Server = ServerConfig{
		Port:         serverPort,
		ReadTimeout:  parseInt(getEnv("READ_TIMEOUT", "30")),
		WriteTimeout: parseInt(getEnv("WRITE_TIMEOUT", "30")),
	}

	config.Debug = strings.ToLower(getEnv("DEBUG", "false")) == "true"

	if config.Database.Password == "" {
		return nil, errors.New("database password is required")
	}

	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return nil, errors.New("server port must be between 1 and 65535")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return result
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
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	LogLevel   string `json:"log_level"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.ServerPort = getEnvOrDefault("SERVER_PORT", config.ServerPort)
	config.DBHost = getEnvOrDefault("DB_HOST", config.DBHost)
	config.DBPort = getEnvOrDefault("DB_PORT", config.DBPort)
	config.LogLevel = getEnvOrDefault("LOG_LEVEL", config.LogLevel)

	return &config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
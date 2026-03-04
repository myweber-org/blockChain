package config

import (
    "fmt"
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
    cfg := &Config{}

    dbHost := getEnv("DB_HOST", "localhost")
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %v", err)
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASS", ""),
        Database: getEnv("DB_NAME", "appdb"),
    }

    serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }

    readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid READ_TIMEOUT: %v", err)
    }

    writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %v", err)
    }

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    debugStr := strings.ToLower(getEnv("DEBUG", "false"))
    cfg.Debug = debugStr == "true" || debugStr == "1"

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("database password is required")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return nil, fmt.Errorf("server port must be between 1 and 65535")
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
	AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var errs []string

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		errs = append(errs, "invalid SERVER_PORT: "+err.Error())
	} else if port < 1 || port > 65535 {
		errs = append(errs, "SERVER_PORT out of range")
	} else {
		cfg.ServerPort = port
	}

	debugStr := os.Getenv("DEBUG_MODE")
	if debugStr == "" {
		debugStr = "false"
	}
	cfg.DebugMode = strings.ToLower(debugStr) == "true"

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}

	ttlStr := os.Getenv("CACHE_TTL")
	if ttlStr == "" {
		ttlStr = "300"
	}
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		errs = append(errs, "invalid CACHE_TTL: "+err.Error())
	} else if ttl < 0 {
		errs = append(errs, "CACHE_TTL cannot be negative")
	} else {
		cfg.CacheTTL = ttl
	}

	hostsStr := os.Getenv("ALLOWED_HOSTS")
	if hostsStr == "" {
		cfg.AllowedHosts = []string{"localhost", "127.0.0.1"}
	} else {
		cfg.AllowedHosts = strings.Split(hostsStr, ",")
		for i, host := range cfg.AllowedHosts {
			cfg.AllowedHosts[i] = strings.TrimSpace(host)
		}
	}

	if len(errs) > 0 {
		return nil, errors.New("configuration errors: " + strings.Join(errs, "; "))
	}

	return cfg, nil
}package config

import (
    "os"
    "strings"
)

type Config struct {
    DatabaseURL string
    APIKey      string
    LogLevel    string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL: getEnvWithDefault("DB_URL", "postgres://localhost:5432/app"),
        APIKey:      getEnvWithDefault("API_KEY", ""),
        LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
    }
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func ParseTemplate(s string) string {
    for _, env := range os.Environ() {
        pair := strings.SplitN(env, "=", 2)
        if len(pair) == 2 {
            placeholder := "${" + pair[0] + "}"
            s = strings.ReplaceAll(s, placeholder, pair[1])
        }
    }
    return s
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
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
    Debug        bool   `yaml:"debug" env:"DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
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
    overrideString(&config.Database.Host, "DB_HOST")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&config.Server.Debug, "DEBUG")
}

func overrideString(field *string, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val
    }
}

func overrideInt(field *int, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        var intVal int
        if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
            *field = intVal
        }
    }
}

func overrideBool(field *bool, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val == "true" || val == "1" || val == "yes"
    }
}

func DefaultConfigPath() string {
    paths := []string{
        "./config.yaml",
        "./config/config.yaml",
        "/etc/app/config.yaml",
    }

    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            absPath, _ := filepath.Abs(path)
            return absPath
        }
    }

    return ""
}package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASS"`
	Database string `yaml:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")
	overrideInt(&config.Server.Port, "SERVER_PORT")
	overrideBool(&config.Server.DebugMode, "DEBUG_MODE")
	overrideString(&config.Database.Host, "DB_HOST")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")
	overrideInt(&config.Database.Port, "DB_PORT")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}

func overrideBool(field *bool, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val == "true" || val == "1" || val == "yes"
	}
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if config.Database.Host == "" {
		return errors.New("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if config.Database.Database == "" {
		return errors.New("database name is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL  string
    MaxConnections int
    DebugMode    bool
    AllowedHosts []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
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
    if value, err := strconv.Atoi(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    if value, err := strconv.ParseBool(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    return strings.Split(strValue, ",")
}package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port    string `yaml:"port" env:"SERVER_PORT"`
		Timeout int    `yaml:"timeout" env:"SERVER_TIMEOUT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(config)
	return config, nil
}

func overrideWithEnvVars(config *Config) {
	overrideStruct(config.Server, "SERVER_")
	overrideStruct(config.Database, "DB_")
	overrideStruct(config.Logging, "LOG_")
}

func overrideStruct(s interface{}, prefix string) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("env")
		if tag == "" {
			continue
		}

		envKey := prefix + tag
		if envValue := os.Getenv(envKey); envValue != "" {
			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)
			case reflect.Int:
				if intVal, err := strconv.Atoi(envValue); err == nil {
					field.SetInt(int64(intVal))
				}
			}
		}
	}
}package config

import (
    "fmt"
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
    cfg := &Config{}

    port, err := getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    cfg.ServerPort = port

    cfg.DBHost = getStringEnv("DB_HOST", "localhost")

    dbPort, err := getIntEnv("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }
    cfg.DBPort = dbPort

    debug, err := getBoolEnv("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }
    cfg.DebugMode = debug

    cfg.APIKeys = getStringSliceEnv("API_KEYS", []string{})

    return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, err
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return false, err
        }
        return boolValue, nil
    }
    return defaultValue, nil
}

func getStringSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}
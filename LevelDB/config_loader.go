package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    ServerPort string `json:"server_port"`
    DatabaseURL string `json:"database_url"`
    LogLevel string `json:"log_level"`
    CacheTTL int `json:"cache_ttl"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config
    
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&cfg); err != nil {
        return nil, err
    }
    
    if port := os.Getenv("SERVER_PORT"); port != "" {
        cfg.ServerPort = port
    }
    
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }
    
    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        cfg.LogLevel = logLevel
    }
    
    return &cfg, nil
}

func SaveConfig(configPath string, cfg *Config) error {
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
    return encoder.Encode(cfg)
}
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DBName     string `env:"DB_NAME" default:"appdb"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		envKey := structField.Tag.Get("env")
		defaultValue := structField.Tag.Get("default")
		
		if envKey == "" {
			continue
		}
		
		envValue := os.Getenv(envKey)
		if envValue == "" {
			envValue = defaultValue
		}
		
		if envValue == "" {
			return nil, fmt.Errorf("environment variable %s is required", envKey)
		}
		
		if err := setFieldValue(field, envValue); err != nil {
			return nil, fmt.Errorf("failed to set %s: %w", envKey, err)
		}
	}
	
	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func (c *Config) Validate() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if c.DBPort < 1 || c.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	
	if !validLogLevels[strings.ToLower(c.LogLevel)] {
		return errors.New("invalid log level")
	}
	
	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
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
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
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
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&config.Server.DebugMode, "DEBUG_MODE")
    overrideString(&config.Server.LogLevel, "LOG_LEVEL")
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
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %v", err)
    }

    return &config, nil
}

func validateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Port)
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    return nil
}
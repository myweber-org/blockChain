package config

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
    Port int    `yaml:"port"`
    Env  string `yaml:"env"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", filePath)
    }

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %v", err)
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    if config.Server.Env == "" {
        config.Server.Env = "development"
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port == 0 {
        return fmt.Errorf("database port is required")
    }
    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}package config

import (
	"errors"
	"os"
	"strconv"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{}
	
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	config.ServerPort = port
	
	config.DBHost = os.Getenv("DB_HOST")
	if config.DBHost == "" {
		config.DBHost = "localhost"
	}
	
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		dbPortStr = "5432"
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}
	config.DBPort = dbPort
	
	debugStr := os.Getenv("DEBUG_MODE")
	if debugStr == "" {
		debugStr = "false"
	}
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		return nil, errors.New("invalid DEBUG_MODE value")
	}
	config.DebugMode = debug
	
	return config, nil
}package config

import (
    "fmt"
    "io/ioutil"
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

func LoadConfig(filePath string) (*ServerConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
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

    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    return nil
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
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if err := overrideFromEnv(&config); err != nil {
        return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
    }

    return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
    overrideString := func(field *string, envVar string) {
        if val := os.Getenv(envVar); val != "" {
            *field = val
        }
    }

    overrideInt := func(field *int, envVar string) error {
        if val := os.Getenv(envVar); val != "" {
            var intVal int
            if _, err := fmt.Sscanf(val, "%d", &intVal); err != nil {
                return fmt.Errorf("invalid integer value for %s: %w", envVar, err)
            }
            *field = intVal
        }
        return nil
    }

    overrideBool := func(field *bool, envVar string) error {
        if val := os.Getenv(envVar); val != "" {
            *field = val == "true" || val == "1" || val == "yes"
        }
        return nil
    }

    db := &config.Database
    if err := overrideInt(&db.Port, "DB_PORT"); err != nil {
        return err
    }
    overrideString(&db.Host, "DB_HOST")
    overrideString(&db.Username, "DB_USER")
    overrideString(&db.Password, "DB_PASS")
    overrideString(&db.Name, "DB_NAME")

    srv := &config.Server
    if err := overrideInt(&srv.Port, "SERVER_PORT"); err != nil {
        return err
    }
    if err := overrideInt(&srv.ReadTimeout, "READ_TIMEOUT"); err != nil {
        return err
    }
    if err := overrideInt(&srv.WriteTimeout, "WRITE_TIMEOUT"); err != nil {
        return err
    }
    if err := overrideBool(&srv.DebugMode, "DEBUG_MODE"); err != nil {
        return err
    }
    overrideString(&srv.LogLevel, "LOG_LEVEL")

    return nil
}

func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", c.Database.Port)
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    if c.Server.ReadTimeout < 0 {
        return fmt.Errorf("read timeout cannot be negative")
    }
    if c.Server.WriteTimeout < 0 {
        return fmt.Errorf("write timeout cannot be negative")
    }
    return nil
}
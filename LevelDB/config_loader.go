package config

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
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    Logging struct {
        Level  string `yaml:"level"`
        Output string `yaml:"output"`
    } `yaml:"logging"`
}

func LoadConfig(filePath string) (*Config, error) {
    config := &Config{}

    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    err = yaml.Unmarshal(data, config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return config, nil
}

func (c *Config) Validate() error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}package config

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

    dbHost := getEnvWithDefault("DB_HOST", "localhost")
    dbPort := getEnvAsInt("DB_PORT", 5432)
    dbUser := getEnvWithDefault("DB_USER", "postgres")
    dbPass := getEnvWithDefault("DB_PASS", "")
    dbName := getEnvWithDefault("DB_NAME", "appdb")

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: dbUser,
        Password: dbPass,
        Database: dbName,
    }

    serverPort := getEnvAsInt("SERVER_PORT", 8080)
    readTimeout := getEnvAsInt("SERVER_READ_TIMEOUT", 30)
    writeTimeout := getEnvAsInt("SERVER_WRITE_TIMEOUT", 30)

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    debugStr := strings.ToLower(getEnvWithDefault("DEBUG", "false"))
    cfg.Debug = debugStr == "true" || debugStr == "1"

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    strValue := getEnvWithDefault(key, "")
    if strValue == "" {
        return defaultValue
    }

    intValue, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return intValue
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
    }

    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }

    if cfg.Server.ReadTimeout <= 0 {
        return fmt.Errorf("read timeout must be positive")
    }

    if cfg.Server.WriteTimeout <= 0 {
        return fmt.Errorf("write timeout must be positive")
    }

    return nil
}package config

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
	Debug    bool           `json:"debug" env:"APP_DEBUG"`
	LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}
	
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	
	overrideFromEnv(config)
	
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	
	return config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideStruct(config, "")
}

func overrideStruct(s interface{}, prefix string) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			if field.Kind() == reflect.Struct {
				overrideStruct(field.Addr().Interface(), prefix)
			}
			continue
		}
		
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		
		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			if intVal, err := strconv.Atoi(envValue); err == nil {
				field.SetInt(int64(intVal))
			}
		case reflect.Bool:
			field.SetBool(strings.ToLower(envValue) == "true")
		}
	}
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("invalid database port")
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	return nil
}package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKey     string
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
	config.DebugMode = strings.ToLower(debugStr) == "true"

	config.APIKey = os.Getenv("API_KEY")
	if config.APIKey == "" {
		return nil, errors.New("API_KEY is required")
	}

	return config, nil
}
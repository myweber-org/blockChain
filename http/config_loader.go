package config

import (
    "fmt"
    "os"
    "strings"

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
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)
    return &config, nil
}

func overrideFromEnv(config *Config) {
    overrideStruct(config.Database)
    overrideStruct(config.Server)
}

func overrideStruct(s interface{}) {
    v := reflect.ValueOf(s).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("env")
        if tag == "" {
            continue
        }

        envValue := os.Getenv(tag)
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
            boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
            field.SetBool(boolVal)
        }
    }
}package config

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Name     string `yaml:"name"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config AppConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        log.Printf("Failed to parse YAML: %v", err)
        return nil, err
    }

    return &config, nil
}
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	Server struct {
		Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
		Port int    `json:"port" env:"SERVER_PORT" default:"8080"`
	} `json:"server"`
	Database struct {
		URL      string `json:"url" env:"DATABASE_URL" required:"true"`
		PoolSize int    `json:"pool_size" env:"DB_POOL_SIZE" default:"10"`
	} `json:"database"`
	LogLevel string `json:"log_level" env:"LOG_LEVEL" default:"info"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}
	
	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}
	
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cfg, nil
}

func loadFromFile(path string, cfg *Config) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}
	
	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}
	
	return nil
}

func loadFromEnv(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		if field.Type.Kind() == reflect.Struct {
			if err := processStruct(fieldValue, field); err != nil {
				return err
			}
		} else {
			if err := processField(fieldValue, field); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func processStruct(structValue reflect.Value, structField reflect.StructField) error {
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		fieldValue := structValue.Field(i)
		
		if err := processField(fieldValue, field); err != nil {
			return err
		}
	}
	return nil
}

func processField(fieldValue reflect.Value, field reflect.StructField) error {
	envTag := field.Tag.Get("env")
	if envTag == "" {
		return nil
	}
	
	envValue := os.Getenv(envTag)
	if envValue == "" {
		defaultTag := field.Tag.Get("default")
		if defaultTag != "" && fieldValue.IsZero() {
			setFieldValue(fieldValue, defaultTag)
		}
		return nil
	}
	
	setFieldValue(fieldValue, envValue)
	return nil
}

func setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal := strings.ToLower(value) == "true" || value == "1"
		field.SetBool(boolVal)
	}
}

func validateConfig(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		if field.Type.Kind() == reflect.Struct {
			if err := validateStruct(fieldValue, field); err != nil {
				return err
			}
		} else {
			if err := validateField(fieldValue, field); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func validateStruct(structValue reflect.Value, structField reflect.StructField) error {
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		fieldValue := structValue.Field(i)
		
		if err := validateField(fieldValue, field); err != nil {
			return err
		}
	}
	return nil
}

func validateField(fieldValue reflect.Value, field reflect.StructField) error {
	required := field.Tag.Get("required")
	if required == "true" && fieldValue.IsZero() {
		envTag := field.Tag.Get("env")
		return fmt.Errorf("field '%s' (env: %s) is required but not set", field.Name, envTag)
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
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int    `yaml:"port"`
    ReadTimeout  int    `yaml:"read_timeout"`
    WriteTimeout int    `yaml:"write_timeout"`
}

type AppConfig struct {
    Environment string         `yaml:"environment"`
    Debug       bool           `yaml:"debug"`
    Database    DatabaseConfig `yaml:"database"`
    Server      ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Environment == "" {
        config.Environment = "development"
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    return &config, nil
}

func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port == 0 {
        return fmt.Errorf("database port is required")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}
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
    SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
    Address      string `json:"address" env:"SERVER_ADDRESS"`
    Port         int    `json:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
}

type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Debug    bool           `json:"debug" env:"DEBUG"`
}

func LoadConfig(configPath string) (*Config, error) {
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, err
    }

    overrideFromEnv(&config)
    
    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func overrideFromEnv(config *Config) {
    overrideStruct(config)
}

func overrideStruct(s interface{}) {
    val := reflect.ValueOf(s).Elem()
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        if field.Kind() == reflect.Struct {
            overrideStruct(field.Addr().Interface())
            continue
        }

        envTag := fieldType.Tag.Get("env")
        if envTag == "" {
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
            boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
            field.SetBool(boolVal)
        }
    }
}

func validateConfig(config *Config) error {
    if config.Database.Host == "" {
        return errors.New("database host is required")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return errors.New("database port must be between 1 and 65535")
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return errors.New("server port must be between 1 and 65535")
    }
    return nil
}package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" validate:"required"`
		Port int    `yaml:"port" validate:"min=1,max=65535"`
	} `yaml:"server"`
	Database struct {
		Driver   string `yaml:"driver" validate:"required"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" validate:"oneof=debug info warn error"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

type Loader struct {
	configPaths []string
	strictMode  bool
}

func NewLoader(paths []string) *Loader {
	return &Loader{
		configPaths: paths,
		strictMode:  true,
	}
}

func (l *Loader) Load() (*Config, error) {
	var configData []byte
	var configFile string

	for _, path := range l.configPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			configData = data
			configFile = path
			break
		}
	}

	if configData == nil {
		return nil, errors.New("no configuration file found in specified paths")
	}

	var cfg Config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from %s: %w", configFile, err)
	}

	if l.strictMode {
		if err := l.validate(&cfg); err != nil {
			return nil, fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	return &cfg, nil
}

func (l *Loader) validate(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return l.validateStruct(v)
}

func (l *Loader) validateStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		if field.Kind() == reflect.Struct {
			if err := l.validateStruct(field); err != nil {
				return err
			}
			continue
		}

		if err := l.validateField(field, tag); err != nil {
			return fmt.Errorf("field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func (l *Loader) validateField(field reflect.Value, rules string) error {
	if rules == "required" {
		if field.IsZero() {
			return errors.New("required field is empty")
		}
	}
	return nil
}

func (l *Loader) SetStrictMode(strict bool) {
	l.strictMode = strict
}

func LoadDefault() (*Config, error) {
	paths := []string{
		"config.yaml",
		"config.yml",
		filepath.Join("config", "config.yaml"),
		filepath.Join("config", "config.yml"),
	}

	loader := NewLoader(paths)
	return loader.Load()
}
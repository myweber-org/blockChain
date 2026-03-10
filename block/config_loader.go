package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Debug    bool           `yaml:"debug" env:"DEBUG"`
}

func LoadConfig(filePath string) (*Config, error) {
    data, err := os.ReadFile(filePath)
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
    overrideStruct(config, "")
}

func overrideStruct(s interface{}, prefix string) {
    v := reflect.ValueOf(s).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        structField := t.Field(i)

        envTag := structField.Tag.Get("env")
        if envTag == "" && field.Kind() == reflect.Struct {
            nestedPrefix := prefix
            if tag := structField.Tag.Get("yaml"); tag != "" {
                nestedPrefix = strings.TrimSuffix(prefix+"_"+strings.ToUpper(tag), "_")
            }
            overrideStruct(field.Addr().Interface(), nestedPrefix)
            continue
        }

        if envTag == "" {
            continue
        }

        envKey := envTag
        if prefix != "" {
            envKey = prefix + "_" + envKey
        }

        if envValue := os.Getenv(envKey); envValue != "" {
            switch field.Kind() {
            case reflect.String:
                field.SetString(envValue)
            case reflect.Int:
                if intVal, err := strconv.Atoi(envValue); err == nil {
                    field.SetInt(int64(intVal))
                }
            case reflect.Bool:
                if boolVal, err := strconv.ParseBool(envValue); err == nil {
                    field.SetBool(boolVal)
                }
            }
        }
    }
}
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST"`
    Port     int    `json:"port" env:"DB_PORT"`
    Username string `json:"username" env:"DB_USER"`
    Password string `json:"password" env:"DB_PASS"`
    Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `json:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
}

type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Debug    bool           `json:"debug" env:"DEBUG"`
}

func LoadConfig(configPath string) (*Config, error) {
    var config Config

    if configPath != "" {
        file, err := os.Open(configPath)
        if err != nil {
            return nil, fmt.Errorf("failed to open config file: %w", err)
        }
        defer file.Close()

        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&config); err != nil {
            return nil, fmt.Errorf("failed to decode config: %w", err)
        }
    }

    loadFromEnv(&config)

    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func loadFromEnv(config *Config) {
    loadStructFromEnv(&config.Database)
    loadStructFromEnv(&config.Server)
    loadStructFromEnv(config)
}

func loadStructFromEnv(s interface{}) {
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
}

func validateConfig(config *Config) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    if config.Server.ReadTimeout < 0 {
        return fmt.Errorf("read timeout cannot be negative")
    }
    if config.Server.WriteTimeout < 0 {
        return fmt.Errorf("write timeout cannot be negative")
    }
    return nil
}
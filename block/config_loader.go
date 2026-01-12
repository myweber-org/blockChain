package config

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
    DebugMode    bool   `yaml:"debug" env:"DEBUG_MODE"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&cfg)

    return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
    overrideString(&cfg.Database.Host, "DB_HOST")
    overrideInt(&cfg.Database.Port, "DB_PORT")
    overrideString(&cfg.Database.Username, "DB_USER")
    overrideString(&cfg.Database.Password, "DB_PASS")
    overrideString(&cfg.Database.Name, "DB_NAME")

    overrideInt(&cfg.Server.Port, "SERVER_PORT")
    overrideInt(&cfg.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&cfg.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&cfg.Server.DebugMode, "DEBUG_MODE")

    overrideString(&cfg.LogLevel, "LOG_LEVEL")
}

func overrideString(field *string, envVar string) {
    if val, exists := os.LookupEnv(envVar); exists && val != "" {
        *field = val
    }
}

func overrideInt(field *int, envVar string) {
    if val, exists := os.LookupEnv(envVar); exists && val != "" {
        var intVal int
        if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
            *field = intVal
        }
    }
}

func overrideBool(field *bool, envVar string) {
    if val, exists := os.LookupEnv(envVar); exists {
        *field = val == "true" || val == "1" || val == "yes"
    }
}
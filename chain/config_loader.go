package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    Database struct {
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(filename string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func (c *AppConfig) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}
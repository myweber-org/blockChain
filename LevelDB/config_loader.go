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
    Port         int            `yaml:"port"`
    Debug        bool           `yaml:"debug"`
    Database     DatabaseConfig `yaml:"database"`
    AllowedHosts []string       `yaml:"allowed_hosts"`
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
        return nil, fmt.Errorf("failed to parse YAML config: %v", err)
    }

    if config.Database.Host == "" {
        config.Database.Host = "localhost"
    }
    if config.Database.Port == 0 {
        config.Database.Port = 5432
    }
    if config.Port == 0 {
        config.Port = 8080
    }

    return &config, nil
}

func (c *ServerConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Port)
    }
    return nil
}
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
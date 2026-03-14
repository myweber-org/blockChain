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
}package config

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
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
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

    if config.Port == 0 {
        config.Port = 8080
    }
    if config.ReadTimeout == 0 {
        config.ReadTimeout = 30
    }
    if config.WriteTimeout == 0 {
        config.WriteTimeout = 30
    }

    return &config, nil
}

func (c *ServerConfig) Validate() error {
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
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
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

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        if port, err := parseInt(val); err == nil {
            config.Database.Port = port
        }
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        if port, err := parseInt(val); err == nil {
            config.Server.Port = port
        }
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        if timeout, err := parseInt(val); err == nil {
            config.Server.ReadTimeout = timeout
        }
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        if timeout, err := parseInt(val); err == nil {
            config.Server.WriteTimeout = timeout
        }
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        config.Server.DebugMode = val == "true" || val == "1"
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }
}

func parseInt(s string) (int, error) {
    var n int
    _, err := fmt.Sscanf(s, "%d", &n)
    return n, err
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
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
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
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
        return nil, fmt.Errorf("failed to apply environment variables: %w", err)
    }

    return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }

    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        var port int
        if _, err := fmt.Sscanf(val, "%d", &port); err == nil {
            config.Database.Port = port
        }
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }

    if val := os.Getenv("SERVER_PORT"); val != "" {
        var port int
        if _, err := fmt.Sscanf(val, "%d", &port); err == nil {
            config.Server.Port = port
        }
    }
    if val := os.Getenv("SERVER_READ_TIMEOUT"); val != "" {
        var timeout int
        if _, err := fmt.Sscanf(val, "%d", &timeout); err == nil {
            config.Server.ReadTimeout = timeout
        }
    }
    if val := os.Getenv("SERVER_WRITE_TIMEOUT"); val != "" {
        var timeout int
        if _, err := fmt.Sscanf(val, "%d", &timeout); err == nil {
            config.Server.WriteTimeout = timeout
        }
    }
    if val := os.Getenv("SERVER_DEBUG"); val != "" {
        config.Server.DebugMode = val == "true" || val == "1"
    }

    return nil
}
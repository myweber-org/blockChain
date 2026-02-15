package config

import (
    "fmt"
    "io/ioutil"
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
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
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
    if c.Database.Port < 1 || c.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
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
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Debug    bool           `yaml:"debug"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    if config.Database.Host == "" {
        config.Database.Host = "localhost"
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Port < 1 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }

    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }

    return nil
}package config

import (
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
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
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    return nil
}package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := make(Config)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax at line %d", lineNumber)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			envVar := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
			envValue := os.Getenv(envVar)
			if envValue == "" {
				return nil, fmt.Errorf("environment variable %s not set at line %d", envVar, lineNumber)
			}
			value = envValue
		}

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

func (c Config) Get(key string) string {
	return c[key]
}

func (c Config) Set(key, value string) {
	c[key] = value
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
    Port        int
    DebugMode   bool
    LogLevel    string
    MaxRequests int
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnvWithDefault("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            Username: getEnvWithDefault("DB_USER", "postgres"),
            Password: getEnvWithDefault("DB_PASSWORD", ""),
            Database: getEnvWithDefault("DB_NAME", "appdb"),
        },
        Server: ServerConfig{
            Port:        getEnvAsInt("SERVER_PORT", 8080),
            DebugMode:   getEnvAsBool("DEBUG_MODE", false),
            LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
            MaxRequests: getEnvAsInt("MAX_REQUESTS", 100),
        },
    }

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
    valueStr := getEnvWithDefault(key, "")
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnvWithDefault(key, "")
    if valueStr == "" {
        return defaultValue
    }
    valueStr = strings.ToLower(valueStr)
    return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }
    if cfg.Server.MaxRequests <= 0 {
        return fmt.Errorf("max requests must be positive: %d", cfg.Server.MaxRequests)
    }
    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    if !validLogLevels[strings.ToLower(cfg.Server.LogLevel)] {
        return fmt.Errorf("invalid log level: %s", cfg.Server.LogLevel)
    }
    return nil
}
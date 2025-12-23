package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
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

func LoadConfig(filePath string) (*Config, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
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
    if c.LogLevel == "" {
        c.LogLevel = "info"
    }
    return nil
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host" env:"SERVER_HOST"`
        Port int    `yaml:"port" env:"SERVER_PORT"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host" env:"DB_HOST"`
        Port     int    `yaml:"port" env:"DB_PORT"`
        Name     string `yaml:"name" env:"DB_NAME"`
        User     string `yaml:"user" env:"DB_USER"`
        Password string `yaml:"password" env:"DB_PASSWORD"`
    } `yaml:"database"`
    Logging struct {
        Level  string `yaml:"level" env:"LOG_LEVEL"`
        Output string `yaml:"output" env:"LOG_OUTPUT"`
    } `yaml:"logging"`
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
    if val := os.Getenv("SERVER_HOST"); val != "" {
        cfg.Server.Host = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.Port)
    }
    if val := os.Getenv("DB_HOST"); val != "" {
        cfg.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Database.Port)
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        cfg.Database.Name = val
    }
    if val := os.Getenv("DB_USER"); val != "" {
        cfg.Database.User = val
    }
    if val := os.Getenv("DB_PASSWORD"); val != "" {
        cfg.Database.Password = val
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        cfg.Logging.Level = val
    }
    if val := os.Getenv("LOG_OUTPUT"); val != "" {
        cfg.Logging.Output = val
    }
}package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	APIKey      string
	DebugMode   bool
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	content = expandEnvVars(content)

	lines := strings.Split(content, "\n")
	cfg := &Config{}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DATABASE_URL":
			cfg.DatabaseURL = value
		case "API_KEY":
			cfg.APIKey = value
		case "DEBUG_MODE":
			cfg.DebugMode = value == "true"
		}
	}

	return cfg, nil
}

func expandEnvVars(s string) string {
	return os.Expand(s, func(key string) string {
		if val, exists := os.LookupEnv(key); exists {
			return val
		}
		return ""
	})
}
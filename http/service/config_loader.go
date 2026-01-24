package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		URL      string `yaml:"url" env:"DB_URL"`
		MaxConns int    `yaml:"max_connections" env:"DB_MAX_CONNS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	
	overrideFromEnv(config)
	
	return config, nil
}

func overrideFromEnv(c *Config) {
	if val := os.Getenv("SERVER_HOST"); val != "" {
		c.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			c.Server.Port = port
		}
	}
	if val := os.Getenv("DB_URL"); val != "" {
		c.Database.URL = val
	}
	if val := os.Getenv("DB_MAX_CONNS"); val != "" {
		if maxConns, err := strconv.Atoi(val); err == nil {
			c.Database.MaxConns = maxConns
		}
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.LogLevel = val
	}
}package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Driver   string `yaml:"driver" env:"DB_DRIVER"`
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		Username string `yaml:"username" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	setFromEnv(&cfg.Server.Host, "SERVER_HOST")
	setFromEnvInt(&cfg.Server.Port, "SERVER_PORT")
	setFromEnv(&cfg.Database.Driver, "DB_DRIVER")
	setFromEnv(&cfg.Database.Host, "DB_HOST")
	setFromEnvInt(&cfg.Database.Port, "DB_PORT")
	setFromEnv(&cfg.Database.Name, "DB_NAME")
	setFromEnv(&cfg.Database.Username, "DB_USERNAME")
	setFromEnv(&cfg.Database.Password, "DB_PASSWORD")
	setFromEnv(&cfg.Logging.Level, "LOG_LEVEL")
	setFromEnv(&cfg.Logging.Output, "LOG_OUTPUT")
}

func setFromEnv(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func setFromEnvInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.ServerPort = getEnvOverride("SERVER_PORT", config.ServerPort)
	config.DBHost = getEnvOverride("DB_HOST", config.DBHost)
	config.DBPort = getEnvOverride("DB_PORT", config.DBPort)

	return &config, nil
}

func getEnvOverride(envKey, defaultValue string) string {
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.ServerPort) == "" {
		return os.ErrInvalid
	}
	if strings.TrimSpace(c.DBHost) == "" {
		return os.ErrInvalid
	}
	return nil
}
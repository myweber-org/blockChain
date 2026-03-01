package config

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
    Port         int
    ReadTimeout  int
    WriteTimeout int
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Debug    bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    dbHost := getEnv("DB_HOST", "localhost")
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %v", err)
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASS", ""),
        Database: getEnv("DB_NAME", "appdb"),
    }

    serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }

    readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid READ_TIMEOUT: %v", err)
    }

    writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %v", err)
    }

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    debugStr := strings.ToLower(getEnv("DEBUG", "false"))
    cfg.Debug = debugStr == "true" || debugStr == "1"

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("database password is required")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return nil, fmt.Errorf("server port must be between 1 and 65535")
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
	AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var errs []string

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		errs = append(errs, "invalid SERVER_PORT: "+err.Error())
	} else if port < 1 || port > 65535 {
		errs = append(errs, "SERVER_PORT out of range")
	} else {
		cfg.ServerPort = port
	}

	debugStr := os.Getenv("DEBUG_MODE")
	if debugStr == "" {
		debugStr = "false"
	}
	cfg.DebugMode = strings.ToLower(debugStr) == "true"

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}

	ttlStr := os.Getenv("CACHE_TTL")
	if ttlStr == "" {
		ttlStr = "300"
	}
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		errs = append(errs, "invalid CACHE_TTL: "+err.Error())
	} else if ttl < 0 {
		errs = append(errs, "CACHE_TTL cannot be negative")
	} else {
		cfg.CacheTTL = ttl
	}

	hostsStr := os.Getenv("ALLOWED_HOSTS")
	if hostsStr == "" {
		cfg.AllowedHosts = []string{"localhost", "127.0.0.1"}
	} else {
		cfg.AllowedHosts = strings.Split(hostsStr, ",")
		for i, host := range cfg.AllowedHosts {
			cfg.AllowedHosts[i] = strings.TrimSpace(host)
		}
	}

	if len(errs) > 0 {
		return nil, errors.New("configuration errors: " + strings.Join(errs, "; "))
	}

	return cfg, nil
}package config

import (
    "os"
    "strings"
)

type Config struct {
    DatabaseURL string
    APIKey      string
    LogLevel    string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL: getEnvWithDefault("DB_URL", "postgres://localhost:5432/app"),
        APIKey:      getEnvWithDefault("API_KEY", ""),
        LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
    }
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func ParseTemplate(s string) string {
    for _, env := range os.Environ() {
        pair := strings.SplitN(env, "=", 2)
        if len(pair) == 2 {
            placeholder := "${" + pair[0] + "}"
            s = strings.ReplaceAll(s, placeholder, pair[1])
        }
    }
    return s
}
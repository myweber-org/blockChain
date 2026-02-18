package config

import (
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  bool   `json:"ssl_mode" env:"DB_SSL"`
}

type AppConfig struct {
	Debug    bool           `json:"debug" env:"APP_DEBUG"`
	LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}
	
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	overrideFromEnv(config)
	
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	
	return config, nil
}

func overrideFromEnv(cfg *AppConfig) {
	overrideStruct(cfg)
}

func overrideStruct(s interface{}) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			overrideStruct(field.Addr().Interface())
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		if envValue := os.Getenv(envTag); envValue != "" {
			setFieldValue(field, envValue)
		}
	}
}

func setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		if intVal, err := strconv.Atoi(value); err == nil {
			field.SetInt(int64(intVal))
		}
	case reflect.Bool:
		boolVal := strings.ToLower(value) == "true" || value == "1"
		field.SetBool(boolVal)
	}
}

func validateConfig(cfg *AppConfig) error {
	if cfg.Database.Host == "" {
		return errors.New("database host is required")
	}
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return errors.New("invalid database port")
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	return nil
}
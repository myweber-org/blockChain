
package main

import (
	"fmt"
	"strings"
)

func FilterAndUppercase(input []string, prefix string) []string {
	var result []string
	for _, s := range input {
		if strings.HasPrefix(s, prefix) {
			result = append(result, strings.ToUpper(s))
		}
	}
	return result
}

func main() {
	data := []string{"apple", "application", "banana", "appetizer"}
	filtered := FilterAndUppercase(data, "app")
	fmt.Println("Processed slice:", filtered)
}
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	Port          int    `json:"port"`
	EnableLogging bool   `json:"enable_logging"`
	MaxConnections int   `json:"max_connections"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(c *Config) error {
	if c.ServerAddress == "" {
		return fmt.Errorf("server_address cannot be empty")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.MaxConnections < 1 {
		return fmt.Errorf("max_connections must be at least 1")
	}
	return nil
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded successfully:\n")
	fmt.Printf("Server: %s:%d\n", config.ServerAddress, config.Port)
	fmt.Printf("Logging enabled: %v\n", config.EnableLogging)
	fmt.Printf("Max connections: %d\n", config.MaxConnections)
}
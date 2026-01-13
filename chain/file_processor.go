
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
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu          sync.Mutex
	processed   int
	errors      []string
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("Failed to open %s: %v", path, err))
		fp.mu.Unlock()
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("Error scanning %s: %v", path, err))
		fp.mu.Unlock()
		return
	}

	fp.mu.Lock()
	fp.processed++
	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	fp.mu.Unlock()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		return
	}

	dir := os.Args[1]
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	processor := &FileProcessor{}
	var wg sync.WaitGroup

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		wg.Add(1)
		go processor.ProcessFile(filePath, &wg)
	}

	wg.Wait()

	fmt.Printf("\nSummary: Processed %d files\n", processor.processed)
	if len(processor.errors) > 0 {
		fmt.Printf("Errors encountered: %d\n", len(processor.errors))
		for _, errMsg := range processor.errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}
}
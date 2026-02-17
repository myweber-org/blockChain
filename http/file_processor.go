package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	Debug  bool   `json:"debug"`
}

func readConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &config, nil
}

func writeConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func main() {
	config := &Config{
		Server: "api.example.com",
		Port:   8080,
		Debug:  true,
	}

	if err := writeConfig("config.json", config); err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config file created successfully")

	loadedConfig, err := readConfig("config.json")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded config: %+v\n", loadedConfig)
}
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int
	Content   string
	Valid     bool
	Timestamp time.Time
}

type Processor struct {
	records []DataRecord
	mu      sync.RWMutex
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(content string) error {
	if content == "" {
		return errors.New("content cannot be empty")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        len(p.records) + 1,
		Content:   content,
		Valid:     true,
		Timestamp: time.Now(),
	}

	p.records = append(p.records, record)
	return nil
}

func (p *Processor) ValidateRecords() {
	var wg sync.WaitGroup
	p.mu.RLock()
	records := make([]DataRecord, len(p.records))
	copy(records, p.records)
	p.mu.RUnlock()

	for i := range records {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p.validateRecord(&records[idx])
		}(i)
	}
	wg.Wait()

	p.mu.Lock()
	p.records = records
	p.mu.Unlock()
}

func (p *Processor) validateRecord(record *DataRecord) {
	if len(record.Content) < 3 {
		record.Valid = false
		return
	}

	time.Sleep(10 * time.Millisecond)
	record.Valid = true
}

func (p *Processor) GetStats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	validCount := 0
	for _, record := range p.records {
		if record.Valid {
			validCount++
		}
	}
	return len(p.records), validCount
}

func main() {
	processor := NewProcessor()

	sampleData := []string{
		"alpha",
		"beta",
		"",
		"gamma",
		"de",
		"epsilon",
	}

	for _, data := range sampleData {
		if err := processor.AddRecord(data); err != nil {
			fmt.Printf("Error adding record: %v\n", err)
		}
	}

	processor.ValidateRecords()
	total, valid := processor.GetStats()
	fmt.Printf("Processing complete. Total: %d, Valid: %d\n", total, valid)
}package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func calculateFileHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		return
	}

	filename := os.Args[1]
	hash, err := calculateFileHash(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	fmt.Printf("SHA256 hash of %s: %s\n", filename, hash)
}package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type DataChunk struct {
	ID    int
	Value string
}

type Processor struct {
	mu     sync.Mutex
	chunks []DataChunk
}

func (p *Processor) AddChunk(chunk DataChunk) error {
	if chunk.ID <= 0 {
		return errors.New("invalid chunk ID")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.chunks = append(p.chunks, chunk)
	return nil
}

func (p *Processor) ProcessAll() {
	var wg sync.WaitGroup
	for _, chunk := range p.chunks {
		wg.Add(1)
		go func(c DataChunk) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Processed chunk %d: %s\n", c.ID, c.Value)
		}(chunk)
	}
	wg.Wait()
}

func main() {
	proc := &Processor{}
	samples := []DataChunk{
		{1, "alpha"},
		{2, "beta"},
		{3, "gamma"},
		{4, "delta"},
	}

	for _, sample := range samples {
		if err := proc.AddChunk(sample); err != nil {
			log.Printf("Failed to add chunk: %v", err)
		}
	}

	proc.ProcessAll()
	fmt.Println("Processing completed")
}
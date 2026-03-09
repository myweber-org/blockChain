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
}package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

type FileProcessor struct {
    inputDir  string
    outputDir string
    workers   int
}

func NewFileProcessor(input, output string, workers int) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
        workers:   workers,
    }
}

func (fp *FileProcessor) ProcessFiles() error {
    files, err := ioutil.ReadDir(fp.inputDir)
    if err != nil {
        return fmt.Errorf("failed to read input directory: %w", err)
    }

    jobs := make(chan string, len(files))
    results := make(chan error, len(files))
    var wg sync.WaitGroup

    for w := 0; w < fp.workers; w++ {
        wg.Add(1)
        go fp.worker(jobs, results, &wg)
    }

    for _, file := range files {
        if !file.IsDir() {
            jobs <- file.Name()
        }
    }
    close(jobs)

    wg.Wait()
    close(results)

    for err := range results {
        if err != nil {
            return err
        }
    }

    return nil
}

func (fp *FileProcessor) worker(jobs <-chan string, results chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()
    for filename := range jobs {
        err := fp.processFile(filename)
        results <- err
    }
}

func (fp *FileProcessor) processFile(filename string) error {
    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, filename)

    data, err := ioutil.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", filename, err)
    }

    processedData := transformData(data)

    if err := os.MkdirAll(fp.outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    if err := ioutil.WriteFile(outputPath, processedData, 0644); err != nil {
        return fmt.Errorf("failed to write file %s: %w", filename, err)
    }

    return nil
}

func transformData(data []byte) []byte {
    result := make([]byte, len(data))
    for i, b := range data {
        result[i] = b ^ 0xFF
    }
    return result
}

func main() {
    processor := NewFileProcessor("./input", "./output", 4)
    if err := processor.ProcessFiles(); err != nil {
        fmt.Printf("Error processing files: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("File processing completed successfully")
}
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu          sync.Mutex
	processed   int
	errors      []error
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		errors: make([]error, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.recordError(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.recordError(err)
		return
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
}

func (fp *FileProcessor) recordError(err error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.errors = append(fp.errors, err)
}

func (fp *FileProcessor) Stats() (int, []error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, fp.errors
}

func ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	processor := NewFileProcessor()
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go processor.ProcessFile(path, &wg)
	}

	wg.Wait()

	processed, errs := processor.Stats()
	fmt.Printf("Total files processed: %d\n", processed)
	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors\n", len(errs))
		for _, err := range errs {
			fmt.Printf("Error: %v\n", err)
		}
		return errors.New("processing completed with errors")
	}

	return nil
}

func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	if err := ProcessFiles(files); err != nil {
		fmt.Fprintf(os.Stderr, "Processing failed: %v\n", err)
		os.Exit(1)
	}
}package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers   int
	batchSize int
	mu        sync.Mutex
	results   map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		results:   make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, len(paths))

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	return nil
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan string) {
	defer wg.Done()

	for file := range files {
		count, err := fp.countLines(file)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", file, err)
			continue
		}

		fp.mu.Lock()
		fp.results[file] = count
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount%fp.batchSize == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func (fp *FileProcessor) GetResults() map[string]int {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	copied := make(map[string]int, len(fp.results))
	for k, v := range fp.results {
		copied[k] = v
	}
	return copied
}

func findTextFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	files, err := findTextFiles(dir)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d text files\n", len(files))

	processor := NewFileProcessor(4, 1000)
	start := time.Now()

	if err := processor.ProcessFiles(files); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	results := processor.GetResults()

	totalLines := 0
	for file, count := range results {
		fmt.Printf("%s: %d lines\n", filepath.Base(file), count)
		totalLines += count
	}

	fmt.Printf("\nTotal lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", elapsed)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	mu          sync.Mutex
	processed   int
	errors      int
	filePattern string
}

func NewFileProcessor(pattern string) *FileProcessor {
	return &FileProcessor{
		filePattern: pattern,
	}
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors++
		fp.mu.Unlock()
		fmt.Printf("Error opening file %s: %v\n", path, err)
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
		fp.errors++
		fp.mu.Unlock()
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", path, lineCount)
}

func (fp *FileProcessor) FindAndProcess() error {
	matches, err := filepath.Glob(fp.filePattern)
	if err != nil {
		return fmt.Errorf("pattern matching error: %v", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("no files found matching pattern: %s", fp.filePattern)
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	for _, match := range matches {
		wg.Add(1)
		go fp.ProcessFile(match, &wg)
	}

	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("Files processed: %d\n", fp.processed)
	fmt.Printf("Errors encountered: %d\n", fp.errors)
	fmt.Printf("Total time: %v\n", duration)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file_pattern>")
		fmt.Println("Example: file_processor *.txt")
		os.Exit(1)
	}

	processor := NewFileProcessor(os.Args[1])
	if err := processor.FindAndProcess(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  map[string]int
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", path, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", path, err)
			return
		}

		fp.mu.Lock()
		fp.results[path] = lineCount
		fp.mu.Unlock()

		fmt.Printf("Processed %s: %d lines\n", path, lineCount)
	}()

	return nil
}

func (fp *FileProcessor) Wait() {
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]int {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			processor.ProcessFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	processor.Wait()

	results := processor.GetResults()
	fmt.Printf("\nTotal files processed: %d\n", len(results))
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "sync"
)

type RecordProcessor interface {
    Process(record []string) error
}

type CSVProcessor struct {
    filePath    string
    concurrency int
    processor   RecordProcessor
}

func NewCSVProcessor(filePath string, concurrency int, processor RecordProcessor) *CSVProcessor {
    return &CSVProcessor{
        filePath:    filePath,
        concurrency: concurrency,
        processor:   processor,
    }
}

func (cp *CSVProcessor) Process() error {
    file, err := os.Open(cp.filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    _, err = reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    var wg sync.WaitGroup
    recordChan := make(chan []string, cp.concurrency*2)
    errorChan := make(chan error, 1)

    for i := 0; i < cp.concurrency; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for record := range recordChan {
                if err := cp.processor.Process(record); err != nil {
                    select {
                    case errorChan <- fmt.Errorf("worker %d: %w", workerID, err):
                    default:
                    }
                    return
                }
            }
        }(i)
    }

    go func() {
        for {
            record, err := reader.Read()
            if err == io.EOF {
                close(recordChan)
                return
            }
            if err != nil {
                select {
                case errorChan <- fmt.Errorf("read error: %w", err):
                default:
                }
                close(recordChan)
                return
            }
            recordChan <- record
        }
    }()

    wg.Wait()
    close(errorChan)

    if err := <-errorChan; err != nil {
        return err
    }
    return nil
}

type StatsProcessor struct {
    mu    sync.Mutex
    count int
}

func (sp *StatsProcessor) Process(record []string) error {
    sp.mu.Lock()
    sp.count++
    sp.mu.Unlock()
    return nil
}

func (sp *StatsProcessor) GetCount() int {
    sp.mu.Lock()
    defer sp.mu.Unlock()
    return sp.count
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: file_processor <csv_file>")
        os.Exit(1)
    }

    processor := &StatsProcessor{}
    csvProcessor := NewCSVProcessor(os.Args[1], 4, processor)

    if err := csvProcessor.Process(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Processed %d records successfully\n", processor.GetCount())
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "sync"
    "time"
)

type DataRecord struct {
    ID        string    `json:"id"`
    Value     float64   `json:"value"`
    Timestamp time.Time `json:"timestamp"`
    Processed bool      `json:"processed"`
}

type Processor struct {
    mu      sync.RWMutex
    records map[string]DataRecord
    logger  *log.Logger
}

func NewProcessor() *Processor {
    return &Processor{
        records: make(map[string]DataRecord),
        logger:  log.New(os.Stdout, "PROCESSOR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

func (p *Processor) AddRecord(id string, value float64) {
    p.mu.Lock()
    defer p.mu.Unlock()

    record := DataRecord{
        ID:        id,
        Value:     value,
        Timestamp: time.Now().UTC(),
        Processed: false,
    }

    p.records[id] = record
    p.logger.Printf("Added record: %s", id)
}

func (p *Processor) ProcessRecord(id string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    record, exists := p.records[id]
    if !exists {
        return fmt.Errorf("record not found: %s", id)
    }

    if record.Processed {
        return fmt.Errorf("record already processed: %s", id)
    }

    record.Processed = true
    record.Value = record.Value * 1.1
    p.records[id] = record

    p.logger.Printf("Processed record: %s", id)
    return nil
}

func (p *Processor) BatchProcess(ids []string) {
    var wg sync.WaitGroup
    errorChan := make(chan error, len(ids))

    for _, id := range ids {
        wg.Add(1)
        go func(recordID string) {
            defer wg.Done()
            if err := p.ProcessRecord(recordID); err != nil {
                errorChan <- err
            }
        }(id)
    }

    wg.Wait()
    close(errorChan)

    for err := range errorChan {
        p.logger.Printf("Processing error: %v", err)
    }
}

func (p *Processor) ExportJSON(filename string) error {
    p.mu.RLock()
    defer p.mu.RUnlock()

    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")

    if err := encoder.Encode(p.records); err != nil {
        return fmt.Errorf("failed to encode JSON: %w", err)
    }

    p.logger.Printf("Exported data to: %s", filename)
    return nil
}

func (p *Processor) Stats() (int, int) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    total := len(p.records)
    processed := 0

    for _, record := range p.records {
        if record.Processed {
            processed++
        }
    }

    return total, processed
}

func main() {
    processor := NewProcessor()

    processor.AddRecord("rec001", 100.0)
    processor.AddRecord("rec002", 200.0)
    processor.AddRecord("rec003", 300.0)

    ids := []string{"rec001", "rec002", "rec003", "rec004"}
    processor.BatchProcess(ids)

    total, processed := processor.Stats()
    fmt.Printf("Total records: %d, Processed: %d\n", total, processed)

    if err := processor.ExportJSON("data_export.json"); err != nil {
        processor.logger.Printf("Export failed: %v", err)
    }
}package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	LogLevel string `json:"log_level"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Server == "" {
		return nil, fmt.Errorf("server address is required")
	}
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("invalid port number: %d", config.Port)
	}
	if config.Timeout < 0 {
		return nil, fmt.Errorf("timeout cannot be negative")
	}

	return &config, nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded successfully:\n")
	fmt.Printf("Server: %s\n", config.Server)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Printf("Timeout: %d seconds\n", config.Timeout)
	fmt.Printf("Log Level: %s\n", config.LogLevel)
}package main

import (
	"fmt"
	"sync"
	"time"
)

type FileProcessor struct {
	mu       sync.Mutex
	wg       sync.WaitGroup
	queue    chan string
	results  map[string]string
}

func NewFileProcessor(workerCount int) *FileProcessor {
	return &FileProcessor{
		queue:   make(chan string, 100),
		results: make(map[string]string),
	}
}

func (fp *FileProcessor) ProcessFile(filename string) string {
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprintf("processed_%s", filename)
}

func (fp *FileProcessor) worker(id int) {
	defer fp.wg.Done()
	for filename := range fp.queue {
		result := fp.ProcessFile(filename)
		
		fp.mu.Lock()
		fp.results[filename] = result
		fp.mu.Unlock()
		
		fmt.Printf("Worker %d processed: %s -> %s\n", id, filename, result)
	}
}

func (fp *FileProcessor) AddFile(filename string) {
	fp.queue <- filename
}

func (fp *FileProcessor) StartWorkers(count int) {
	fp.wg.Add(count)
	for i := 1; i <= count; i++ {
		go fp.worker(i)
	}
}

func (fp *FileProcessor) Wait() {
	close(fp.queue)
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]string {
	return fp.results
}

func main() {
	processor := NewFileProcessor(3)
	processor.StartWorkers(3)
	
	files := []string{"document1.txt", "image2.jpg", "data3.csv", "report4.pdf"}
	
	for _, file := range files {
		processor.AddFile(file)
	}
	
	processor.Wait()
	
	fmt.Println("\nProcessing complete. Results:")
	for filename, result := range processor.GetResults() {
		fmt.Printf("%s: %s\n", filename, result)
	}
}
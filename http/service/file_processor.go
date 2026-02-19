package main

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
	mu          sync.RWMutex
	records     map[string]DataRecord
	workerCount int
}

func NewProcessor(workers int) *Processor {
	return &Processor{
		records:     make(map[string]DataRecord),
		workerCount: workers,
	}
}

func (p *Processor) LoadData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var records []DataRecord
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	p.mu.Lock()
	for _, record := range records {
		p.records[record.ID] = record
	}
	p.mu.Unlock()

	log.Printf("Loaded %d records from %s", len(records), filename)
	return nil
}

func (p *Processor) ProcessRecord(id string) error {
	p.mu.RLock()
	record, exists := p.records[id]
	p.mu.RUnlock()

	if !exists {
		return fmt.Errorf("record with ID %s not found", id)
	}

	if record.Processed {
		return fmt.Errorf("record %s already processed", id)
	}

	time.Sleep(50 * time.Millisecond)

	record.Value = record.Value * 1.1
	record.Processed = true

	p.mu.Lock()
	p.records[id] = record
	p.mu.Unlock()

	return nil
}

func (p *Processor) RunConcurrentProcessing() {
	var wg sync.WaitGroup
	ids := p.getAllIDs()

	workChan := make(chan string, len(ids))

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for id := range workChan {
				if err := p.ProcessRecord(id); err != nil {
					log.Printf("Worker %d: Error processing %s: %v", workerID, id, err)
				} else {
					log.Printf("Worker %d: Successfully processed %s", workerID, id)
				}
			}
		}(i)
	}

	for _, id := range ids {
		workChan <- id
	}
	close(workChan)

	wg.Wait()
	log.Println("All processing completed")
}

func (p *Processor) getAllIDs() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ids := make([]string, 0, len(p.records))
	for id := range p.records {
		ids = append(ids, id)
	}
	return ids
}

func (p *Processor) SaveResults(filename string) error {
	p.mu.RLock()
	records := make([]DataRecord, 0, len(p.records))
	for _, record := range p.records {
		records = append(records, record)
	}
	p.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	log.Printf("Saved %d records to %s", len(records), filename)
	return nil
}

func generateSampleData() []DataRecord {
	records := make([]DataRecord, 100)
	for i := range records {
		records[i] = DataRecord{
			ID:        fmt.Sprintf("REC-%04d", i+1),
			Value:     float64(i+1) * 10.0,
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
			Processed: false,
		}
	}
	return records
}

func main() {
	sampleFile := "sample_data.json"
	outputFile := "processed_data.json"

	sampleData := generateSampleData()
	file, err := os.Create(sampleFile)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(sampleData)
	file.Close()

	processor := NewProcessor(5)

	if err := processor.LoadData(sampleFile); err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	processor.RunConcurrentProcessing()
	duration := time.Since(startTime)

	log.Printf("Processing took %v", duration)

	if err := processor.SaveResults(outputFile); err != nil {
		log.Fatal(err)
	}

	os.Remove(sampleFile)
	log.Println("Sample data cleaned up")
}
package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type FileResult struct {
    Path    string
    Size    int64
    Lines   int
    Error   error
}

func processFile(path string, results chan<- FileResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    result := FileResult{Path: path}
    
    file, err := os.Open(path)
    if err != nil {
        result.Error = err
        results <- result
        return
    }
    defer file.Close()
    
    info, err := file.Stat()
    if err != nil {
        result.Error = err
        results <- result
        return
    }
    result.Size = info.Size()
    
    scanner := bufio.NewScanner(file)
    lineCount := 0
    for scanner.Scan() {
        lineCount++
    }
    
    if err := scanner.Err(); err != nil {
        result.Error = err
        results <- result
        return
    }
    
    result.Lines = lineCount
    results <- result
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: file_processor <directory>")
        os.Exit(1)
    }
    
    root := os.Args[1]
    var files []string
    
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    
    if err != nil {
        fmt.Printf("Error walking directory: %v\n", err)
        os.Exit(1)
    }
    
    results := make(chan FileResult, len(files))
    var wg sync.WaitGroup
    
    for _, file := range files {
        wg.Add(1)
        go processFile(file, results, &wg)
    }
    
    wg.Wait()
    close(results)
    
    totalFiles := 0
    totalSize := int64(0)
    totalLines := 0
    
    for result := range results {
        if result.Error != nil {
            fmt.Printf("Error processing %s: %v\n", result.Path, result.Error)
            continue
        }
        
        totalFiles++
        totalSize += result.Size
        totalLines += result.Lines
        
        fmt.Printf("%s: %d bytes, %d lines\n", 
            filepath.Base(result.Path), result.Size, result.Lines)
    }
    
    fmt.Printf("\nSummary: %d files, %d total bytes, %d total lines\n", 
        totalFiles, totalSize, totalLines)
}
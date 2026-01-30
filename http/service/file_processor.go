
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	Workers    int
	BatchSize  int
	Results    chan string
	Errors     chan error
	wg         sync.WaitGroup
	mu         sync.Mutex
	processed  int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:   workers,
		BatchSize: batchSize,
		Results:   make(chan string, 100),
		Errors:    make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	fileChan := make(chan string, len(files))
	for _, file := range files {
		if !file.IsDir() {
			fileChan <- filepath.Join(dirPath, file.Name())
		}
	}
	close(fileChan)

	for i := 0; i < fp.Workers; i++ {
		fp.wg.Add(1)
		go fp.worker(fileChan)
	}

	fp.wg.Wait()
	close(fp.Results)
	close(fp.Errors)

	return nil
}

func (fp *FileProcessor) worker(files <-chan string) {
	defer fp.wg.Done()

	batch := make([]string, 0, fp.BatchSize)
	for file := range files {
		batch = append(batch, file)

		if len(batch) >= fp.BatchSize {
			fp.processBatch(batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		fp.processBatch(batch)
	}
}

func (fp *FileProcessor) processBatch(files []string) {
	var batchWg sync.WaitGroup
	batchWg.Add(len(files))

	for _, file := range files {
		go func(f string) {
			defer batchWg.Done()
			if err := fp.processSingleFile(f); err != nil {
				fp.Errors <- fmt.Errorf("file %s: %w", f, err)
			}
		}(file)
	}

	batchWg.Wait()

	fp.mu.Lock()
	fp.processed += len(files)
	fmt.Printf("Processed batch of %d files, total: %d\n", len(files), fp.processed)
	fp.mu.Unlock()
}

func (fp *FileProcessor) processSingleFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		_ = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if lineCount == 0 {
		return errors.New("empty file")
	}

	fp.Results <- fmt.Sprintf("%s: %d lines", filepath.Base(filePath), lineCount)
	return nil
}

func (fp *FileProcessor) Stats() {
	fmt.Printf("Total files processed: %d\n", fp.processed)
}

func main() {
	processor := NewFileProcessor(4, 10)

	go func() {
		for result := range processor.Results {
			fmt.Println("Result:", result)
		}
	}()

	go func() {
		for err := range processor.Errors {
			fmt.Println("Error:", err)
		}
	}()

	start := time.Now()
	if err := processor.ProcessDirectory("."); err != nil {
		fmt.Println("Processing error:", err)
		return
	}
	elapsed := time.Since(start)

	processor.Stats()
	fmt.Printf("Processing completed in %v\n", elapsed)
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "sync"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

func processFile(filename string, results chan<- Record, errors chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()

    file, err := os.Open(filename)
    if err != nil {
        errors <- fmt.Errorf("failed to open file: %w", err)
        return
    }
    defer file.Close()

    reader := csv.NewReader(file)
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            errors <- fmt.Errorf("line %d: csv read error: %w", lineNumber, err)
            continue
        }

        if len(row) != 3 {
            errors <- fmt.Errorf("line %d: invalid column count %d", lineNumber, len(row))
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            errors <- fmt.Errorf("line %d: invalid ID format: %w", lineNumber, err)
            continue
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            errors <- fmt.Errorf("line %d: invalid value format: %w", lineNumber, err)
            continue
        }

        results <- Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        }
    }
}

func aggregateResults(results <-chan Record, errors <-chan error) ([]Record, []error) {
    var records []Record
    var errs []error

    for {
        select {
        case record, ok := <-results:
            if !ok {
                results = nil
            } else {
                records = append(records, record)
            }
        case err, ok := <-errors:
            if !ok {
                errors = nil
            } else {
                errs = append(errs, err)
            }
        }

        if results == nil && errors == nil {
            break
        }
    }

    return records, errs
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: file_processor <filename>")
        os.Exit(1)
    }

    filename := os.Args[1]
    results := make(chan Record, 100)
    errors := make(chan error, 100)
    var wg sync.WaitGroup

    wg.Add(1)
    go processFile(filename, results, errors, &wg)

    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    records, errs := aggregateResults(results, errors)

    fmt.Printf("Processed %d records\n", len(records))
    if len(errs) > 0 {
        fmt.Printf("Encountered %d errors:\n", len(errs))
        for _, err := range errs {
            fmt.Printf("  - %v\n", err)
        }
    }

    if len(records) > 0 {
        fmt.Println("\nFirst 5 records:")
        for i := 0; i < len(records) && i < 5; i++ {
            fmt.Printf("  ID: %d, Name: %s, Value: %.2f\n",
                records[i].ID, records[i].Name, records[i].Value)
        }
    }
}
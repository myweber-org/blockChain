
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
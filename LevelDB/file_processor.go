package main

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
	stats     map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(i, fileChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	fp.printStats()
	return nil
}

func (fp *FileProcessor) worker(id int, files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		start := time.Now()
		err := fp.processSingleFile(file)
		duration := time.Since(start)

		fp.mu.Lock()
		if err != nil {
			fp.stats["errors"]++
			fmt.Printf("Worker %d: Failed to process %s: %v\n", id, file, err)
		} else {
			fp.stats["processed"]++
			fmt.Printf("Worker %d: Processed %s in %v\n", id, file, duration)
		}
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) processSingleFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	switch ext {
	case ".txt":
		return fp.processTextFile(file)
	case ".log":
		return fp.processLogFile(file)
	default:
		return fp.processGenericFile(file)
	}
}

func (fp *FileProcessor) processTextFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		_ = scanner.Text()
	}

	fp.mu.Lock()
	fp.stats["lines"] += lineCount
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) processLogFile(file *os.File) error {
	_, err := io.Copy(io.Discard, file)
	return err
}

func (fp *FileProcessor) processGenericFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	fp.mu.Lock()
	fp.stats["bytes"] += int(info.Size())
	fp.mu.Unlock()

	return nil
}

func (fp *FileProcessor) printStats() {
	fmt.Println("\nProcessing Statistics:")
	fmt.Println("=====================")
	for key, value := range fp.stats {
		fmt.Printf("%s: %d\n", key, value)
	}
}

func main() {
	processor := NewFileProcessor(4, 10)

	files := []string{
		"data/file1.txt",
		"data/file2.log",
		"data/file3.csv",
		"data/file4.txt",
	}

	if err := processor.ProcessFiles(files); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}
}
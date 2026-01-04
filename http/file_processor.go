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

type FileStats struct {
	Path        string
	Size        int64
	LineCount   int
	ProcessTime time.Duration
	Error       error
}

type FileProcessor struct {
	workers    int
	results    chan FileStats
	wg         sync.WaitGroup
	mu         sync.Mutex
	totalFiles int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		results: make(chan FileStats, 100),
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dirPath)
	}

	fp.wg.Add(fp.workers)
	for i := 0; i < fp.workers; i++ {
		go fp.worker()
	}

	fileChan := make(chan string, 1000)
	go fp.walkDirectory(dirPath, fileChan)

	go func() {
		fp.wg.Wait()
		close(fp.results)
	}()

	var totalSize int64
	var totalLines int
	var failedFiles int

	for stats := range fp.results {
		if stats.Error != nil {
			failedFiles++
			fmt.Printf("Error processing %s: %v\n", stats.Path, stats.Error)
			continue
		}
		totalSize += stats.Size
		totalLines += stats.LineCount
		fmt.Printf("Processed %s: %d bytes, %d lines in %v\n",
			stats.Path, stats.Size, stats.LineCount, stats.ProcessTime)
	}

	fmt.Printf("\nSummary: Processed %d files, %d failed\n", fp.totalFiles, failedFiles)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)

	return nil
}

func (fp *FileProcessor) walkDirectory(dirPath string, fileChan chan<- string) {
	defer close(fileChan)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		fileChan <- path
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}

func (fp *FileProcessor) worker() {
	defer fp.wg.Done()

	for path := range fp.results {
		stats := fp.processFile(path)
		fp.results <- stats
	}
}

func (fp *FileProcessor) processFile(path string) FileStats {
	start := time.Now()
	stats := FileStats{Path: path}

	file, err := os.Open(path)
	if err != nil {
		stats.Error = err
		stats.ProcessTime = time.Since(start)
		return stats
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		stats.Error = err
		stats.ProcessTime = time.Since(start)
		return stats
	}
	stats.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount%1000 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		stats.Error = err
		stats.ProcessTime = time.Since(start)
		return stats
	}

	stats.LineCount = lineCount
	stats.ProcessTime = time.Since(start)

	fp.mu.Lock()
	fp.totalFiles++
	fp.mu.Unlock()

	return stats
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor(4)
	if err := processor.ProcessDirectory(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
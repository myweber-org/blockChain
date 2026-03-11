package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	fileList []string
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		fileList: make([]string, 0),
	}
}

func (fp *FileProcessor) ScanDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fp.mu.Lock()
			fp.fileList = append(fp.fileList, path)
			fp.mu.Unlock()
		}
		return nil
	})
}

func (fp *FileProcessor) ProcessFiles() []string {
	var wg sync.WaitGroup
	results := make([]string, 0)
	resultChan := make(chan string, len(fp.fileList))

	for _, file := range fp.fileList {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			processed := fp.processSingleFile(f)
			resultChan <- processed
		}(file)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func (fp *FileProcessor) processSingleFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("SCAN_ERROR: %s", err.Error())
	}

	return fmt.Sprintf("File: %s, Lines: %d", filepath.Base(filePath), lineCount)
}

func (fp *FileProcessor) GetFileCount() int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return len(fp.fileList)
}

func main() {
	processor := NewFileProcessor()
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		return
	}

	dirPath := os.Args[1]
	
	err := processor.ScanDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	fmt.Printf("Found %d files\n", processor.GetFileCount())
	
	results := processor.ProcessFiles()
	
	for _, result := range results {
		fmt.Println(result)
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
	results  []string
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make([]string, 0),
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

		fp.mu.Lock()
		fp.results = append(fp.results, fmt.Sprintf("%s: %d lines", filepath.Base(path), lineCount))
		fp.mu.Unlock()
	}()

	return nil
}

func (fp *FileProcessor) Wait() []string {
	fp.wg.Wait()
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> <file2> ...")
		os.Exit(1)
	}

	processor := NewFileProcessor()

	for _, filePath := range os.Args[1:] {
		if err := processor.ProcessFile(filePath); err != nil {
			fmt.Printf("Failed to process %s: %v\n", filePath, err)
		}
	}

	results := processor.Wait()
	fmt.Println("Processing results:")
	for _, result := range results {
		fmt.Println(result)
	}
}package main

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
	Workers   int
	BatchSize int
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 3
	}
	if batchSize < 1 {
		batchSize = 10
	}
	return &FileProcessor{
		Workers:   workers,
		BatchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, action func(string) error) []error {
	fileChan := make(chan string, len(paths))
	errChan := make(chan error, len(paths))
	resultChan := make(chan []error, fp.Workers)

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	for i := 0; i < fp.Workers; i++ {
		fp.wg.Add(1)
		go fp.worker(fileChan, errChan, action)
	}

	go func() {
		fp.wg.Wait()
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (fp *FileProcessor) worker(files <-chan string, errChan chan<- error, action func(string) error) {
	defer fp.wg.Done()
	
	batch := make([]string, 0, fp.BatchSize)
	
	for file := range files {
		batch = append(batch, file)
		
		if len(batch) >= fp.BatchSize {
			fp.processBatch(batch, errChan, action)
			batch = batch[:0]
		}
	}
	
	if len(batch) > 0 {
		fp.processBatch(batch, errChan, action)
	}
}

func (fp *FileProcessor) processBatch(batch []string, errChan chan<- error, action func(string) error) {
	var batchWg sync.WaitGroup
	batchWg.Add(len(batch))
	
	for _, file := range batch {
		go func(f string) {
			defer batchWg.Done()
			if err := action(f); err != nil {
				errChan <- fmt.Errorf("process %s: %w", f, err)
			}
		}(file)
	}
	batchWg.Wait()
}

func ValidateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file does not exist")
		}
		return err
	}
	
	if info.IsDir() {
		return errors.New("path is a directory")
	}
	
	if info.Size() == 0 {
		return errors.New("file is empty")
	}
	
	return nil
}

func CountLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineCount := 0
	
	for scanner.Scan() {
		lineCount++
	}
	
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	
	return lineCount, nil
}

func main() {
	processor := NewFileProcessor(4, 5)
	
	files := []string{
		"data1.txt",
		"data2.txt",
		"data3.txt",
		"data4.txt",
		"data5.txt",
	}
	
	start := time.Now()
	
	errors := processor.ProcessFiles(files, func(path string) error {
		if err := ValidateFile(path); err != nil {
			return err
		}
		
		lines, err := CountLines(path)
		if err != nil {
			return err
		}
		
		absPath, _ := filepath.Abs(path)
		fmt.Printf("File: %s, Lines: %d\n", absPath, lines)
		return nil
	})
	
	elapsed := time.Since(start)
	
	if len(errors) > 0 {
		fmt.Printf("Completed with %d errors in %v:\n", len(errors), elapsed)
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	} else {
		fmt.Printf("All files processed successfully in %v\n", elapsed)
	}
}
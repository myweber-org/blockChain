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
	mu          sync.Mutex
	processed   int
	errors      []error
	concurrency int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		concurrency: workers,
		errors:      make([]error, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning %s: %w", path, err)
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	return nil
}

func (fp *FileProcessor) ProcessDirectory(dir string) error {
	var wg sync.WaitGroup
	fileChan := make(chan string, fp.concurrency)

	for i := 0; i < fp.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				if err := fp.ProcessFile(path); err != nil {
					fp.mu.Lock()
					fp.errors = append(fp.errors, err)
					fp.mu.Unlock()
				}
			}
		}()
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		fileChan <- path
		return nil
	})

	close(fileChan)
	wg.Wait()

	if err != nil {
		return fmt.Errorf("directory walk error: %w", err)
	}

	if len(fp.errors) > 0 {
		return errors.New("some files failed to process")
	}

	return nil
}

func (fp *FileProcessor) Stats() (int, int) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, len(fp.errors)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	start := time.Now()
	processor := NewFileProcessor(4)

	dir := os.Args[1]
	fmt.Printf("Processing files in: %s\n", dir)

	if err := processor.ProcessDirectory(dir); err != nil {
		fmt.Printf("Processing completed with errors: %v\n", err)
	}

	processed, errors := processor.Stats()
	elapsed := time.Since(start)

	fmt.Printf("\nResults:\n")
	fmt.Printf("  Files processed: %d\n", processed)
	fmt.Printf("  Errors: %d\n", errors)
	fmt.Printf("  Time elapsed: %v\n", elapsed)
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	InputDir  string
	OutputDir string
	Workers   int
}

func NewFileProcessor(inputDir, outputDir string, workers int) *FileProcessor {
	return &FileProcessor{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Workers:   workers,
	}
}

func (fp *FileProcessor) ProcessFiles() error {
	files, err := filepath.Glob(filepath.Join(fp.InputDir, "*.txt"))
	if err != nil {
		return err
	}

	fileChan := make(chan string, len(files))
	resultChan := make(chan string, len(files))
	var wg sync.WaitGroup

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(i, fileChan, resultChan, &wg)
	}

	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println("Processed:", result)
	}

	return nil
}

func (fp *FileProcessor) worker(id int, files <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		start := time.Now()
		outputFile := filepath.Join(fp.OutputDir, fmt.Sprintf("processed_%d_%s", id, filepath.Base(file)))

		if err := fp.processSingleFile(file, outputFile); err != nil {
			results <- fmt.Sprintf("Worker %d failed on %s: %v", id, file, err)
			continue
		}

		duration := time.Since(start)
		results <- fmt.Sprintf("Worker %d completed %s in %v", id, file, duration)
	}
}

func (fp *FileProcessor) processSingleFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		processedLine := fmt.Sprintf("PROCESSED: %s\n", line)
		if _, err := writer.WriteString(processedLine); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}

func main() {
	processor := NewFileProcessor("./input", "./output", 4)
	if err := processor.ProcessFiles(); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All files processed successfully")
}
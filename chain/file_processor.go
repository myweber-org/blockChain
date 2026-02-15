package main

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

    var wg sync.WaitGroup
    fileChan := make(chan os.FileInfo, len(files))
    errChan := make(chan error, fp.workers)

    for i := 0; i < fp.workers; i++ {
        wg.Add(1)
        go fp.worker(&wg, fileChan, errChan)
    }

    for _, file := range files {
        if !file.IsDir() {
            fileChan <- file
        }
    }
    close(fileChan)

    wg.Wait()
    close(errChan)

    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan os.FileInfo, errChan chan<- error) {
    defer wg.Done()

    for file := range files {
        if err := fp.processSingleFile(file.Name()); err != nil {
            errChan <- err
        }
    }
}

func (fp *FileProcessor) processSingleFile(filename string) error {
    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, filename)

    data, err := ioutil.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", filename, err)
    }

    processedData := transformData(data)

    if err := ioutil.WriteFile(outputPath, processedData, 0644); err != nil {
        return fmt.Errorf("failed to write file %s: %w", filename, err)
    }

    fmt.Printf("Processed: %s -> %s\n", inputPath, outputPath)
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
	wg        sync.WaitGroup
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, handler func(string) error) []error {
	var errs []error
	jobChan := make(chan string, len(paths))
	resultChan := make(chan error, len(paths))

	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(jobChan, resultChan, handler)
	}

	for _, path := range paths {
		jobChan <- path
	}
	close(jobChan)

	fp.wg.Wait()
	close(resultChan)

	for err := range resultChan {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (fp *FileProcessor) worker(jobs <-chan string, results chan<- error, handler func(string) error) {
	defer fp.wg.Done()
	for path := range jobs {
		results <- handler(path)
	}
}

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return lines, nil
}

func validateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.New("file does not exist")
	}
	if info.IsDir() {
		return errors.New("path is a directory")
	}
	if info.Size() == 0 {
		return errors.New("file is empty")
	}
	return nil
}

func processSingleFile(path string) error {
	if err := validateFile(path); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	lines, err := readFileLines(path)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	fp.mu.Lock()
	defer fp.mu.Unlock()
	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), len(lines))
	return nil
}

func main() {
	processor := NewFileProcessor(4, 10)
	files := []string{"data1.txt", "data2.txt", "data3.txt"}

	start := time.Now()
	errors := processor.ProcessFiles(files, processSingleFile)
	elapsed := time.Since(start)

	if len(errors) > 0 {
		fmt.Printf("Completed with %d errors in %v\n", len(errors), elapsed)
		for _, err := range errors {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Printf("All files processed successfully in %v\n", elapsed)
	}
}package main

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
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &config, nil
}

func writeConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
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

	loadedConfig, err := readConfig("config.json")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded config: %+v\n", loadedConfig)
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

type FileStats struct {
	LineCount int
	WordCount int
	CharCount int
	FilePath  string
	Processed bool
	Error     error
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	stats := FileStats{FilePath: path}

	file, err := os.Open(path)
	if err != nil {
		stats.Error = err
		results <- stats
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineScanner := bufio.NewScanner(reader)

	for lineScanner.Scan() {
		line := lineScanner.Text()
		stats.LineCount++
		stats.CharCount += len(line)

		wordScanner := bufio.NewScanner(bufio.NewReader(filepath.Clean(line)))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			stats.WordCount++
		}
	}

	if err := lineScanner.Err(); err != nil {
		stats.Error = err
		results <- stats
		return
	}

	stats.Processed = true
	results <- stats
}

func aggregateResults(results <-chan FileStats, totalFiles int) ([]FileStats, error) {
	processedFiles := make([]FileStats, 0, totalFiles)
	var finalErr error

	for i := 0; i < totalFiles; i++ {
		stats := <-results
		processedFiles = append(processedFiles, stats)

		if stats.Error != nil && finalErr == nil {
			finalErr = stats.Error
		}
	}

	if finalErr != nil {
		return processedFiles, finalErr
	}
	return processedFiles, nil
}

func validatePaths(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no file paths provided")
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		os.Exit(1)
	}

	filePaths := os.Args[1:]

	if err := validatePaths(filePaths); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	results := make(chan FileStats, len(filePaths))

	for _, path := range filePaths {
		wg.Add(1)
		go processFile(path, &wg, results)
	}

	wg.Wait()
	close(results)

	processedStats, err := aggregateResults(results, len(filePaths))
	if err != nil {
		fmt.Printf("Processing completed with errors: %v\n", err)
	}

	totalLines := 0
	totalWords := 0
	totalChars := 0
	successCount := 0

	for _, stats := range processedStats {
		if stats.Processed {
			successCount++
			totalLines += stats.LineCount
			totalWords += stats.WordCount
			totalChars += stats.CharCount

			fmt.Printf("File: %s\n", stats.FilePath)
			fmt.Printf("  Lines: %d, Words: %d, Characters: %d\n",
				stats.LineCount, stats.WordCount, stats.CharCount)
		} else if stats.Error != nil {
			fmt.Printf("File: %s - Error: %v\n", stats.FilePath, stats.Error)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Files processed successfully: %d/%d\n", successCount, len(filePaths))
	fmt.Printf("  Total lines: %d\n", totalLines)
	fmt.Printf("  Total words: %d\n", totalWords)
	fmt.Printf("  Total characters: %d\n", totalChars)
	fmt.Printf("  Processing time: %v\n", duration)
}
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
}
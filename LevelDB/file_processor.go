
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

func (fp *FileProcessor) ProcessFiles() {
	var wg sync.WaitGroup
	results := make(chan string, len(fp.fileList))

	for _, file := range fp.fileList {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			result := fp.processSingleFile(f)
			results <- result
		}(file)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println(result)
	}
}

func (fp *FileProcessor) processSingleFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("Error reading file %s: %v", filePath, err)
	}

	return fmt.Sprintf("File: %s, Lines: %d", filepath.Base(filePath), lineCount)
}

func (fp *FileProcessor) GetFileCount() int {
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
	processor.ProcessFiles()
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
		
		fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
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
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		os.Exit(1)
	}
	
	processor := NewFileProcessor()
	
	for _, filePath := range os.Args[1:] {
		if err := processor.ProcessFile(filePath); err != nil {
			fmt.Printf("Failed to process %s: %v\n", filePath, err)
		}
	}
	
	processor.Wait()
	
	fmt.Println("\nProcessing complete. Results:")
	for file, lines := range processor.GetResults() {
		fmt.Printf("%s: %d lines\n", filepath.Base(file), lines)
	}
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
	BatchSize   int
	Workers     int
	mu          sync.Mutex
	processed   int
	errors      []error
}

func NewFileProcessor(batchSize, workers int) *FileProcessor {
	if batchSize <= 0 {
		batchSize = 100
	}
	if workers <= 0 {
		workers = 4
	}
	return &FileProcessor{
		BatchSize: batchSize,
		Workers:   workers,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, len(paths))
	resultChan := make(chan error, len(paths))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan, resultChan, processor)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	close(resultChan)

	for err := range resultChan {
		if err != nil {
			fp.mu.Lock()
			fp.errors = append(fp.errors, err)
			fp.mu.Unlock()
		}
	}

	if len(fp.errors) > 0 {
		return fmt.Errorf("encountered %d errors during processing", len(fp.errors))
	}
	return nil
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan string, results chan<- error, processor func(string) error) {
	defer wg.Done()
	for file := range files {
		err := processor(file)
		results <- err
		fp.mu.Lock()
		fp.processed++
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

func (fp *FileProcessor) CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (fp *FileProcessor) FindFilesByPattern(dir, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func (fp *FileProcessor) Stats() (processed int, errorCount int) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, len(fp.errors)
}

func main() {
	processor := NewFileProcessor(50, 8)
	
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	
	err := processor.ProcessFiles(files, func(path string) error {
		lines, err := processor.CountLines(path)
		if err != nil {
			return err
		}
		fmt.Printf("File %s has %d lines\n", path, lines)
		return nil
	})
	
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
	}
	
	processed, errors := processor.Stats()
	fmt.Printf("Processed %d files with %d errors\n", processed, errors)
	
	time.Sleep(100 * time.Millisecond)
}
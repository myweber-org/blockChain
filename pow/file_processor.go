
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	mu          sync.Mutex
	processed   int
	errors      int
	filePattern string
}

func NewFileProcessor(pattern string) *FileProcessor {
	return &FileProcessor{
		filePattern: pattern,
	}
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors++
		fp.mu.Unlock()
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
		fp.mu.Lock()
		fp.errors++
		fp.mu.Unlock()
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", path, lineCount)
}

func (fp *FileProcessor) FindAndProcess() error {
	matches, err := filepath.Glob(fp.filePattern)
	if err != nil {
		return fmt.Errorf("pattern matching error: %v", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("no files found matching pattern: %s", fp.filePattern)
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	for _, match := range matches {
		wg.Add(1)
		go fp.ProcessFile(match, &wg)
	}

	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("Files processed: %d\n", fp.processed)
	fmt.Printf("Errors encountered: %d\n", fp.errors)
	fmt.Printf("Total time: %v\n", duration)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file_pattern>")
		fmt.Println("Example: file_processor *.txt")
		os.Exit(1)
	}

	processor := NewFileProcessor(os.Args[1])
	if err := processor.FindAndProcess(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
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

		fmt.Printf("Processed %s: %d lines\n", path, lineCount)
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
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			processor.ProcessFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	processor.Wait()

	results := processor.GetResults()
	fmt.Printf("\nTotal files processed: %d\n", len(results))
}

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
		
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", path, err)
			return
		}
		
		result := fmt.Sprintf("File: %s, Lines: %d", filepath.Base(path), lineCount)
		
		fp.mu.Lock()
		fp.results = append(fp.results, result)
		fp.mu.Unlock()
	}()
	
	return nil
}

func (fp *FileProcessor) Wait() {
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() []string {
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
	
	processor.Wait()
	
	results := processor.GetResults()
	fmt.Println("Processing complete:")
	for _, result := range results {
		fmt.Println(result)
	}
}
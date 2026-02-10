
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
}
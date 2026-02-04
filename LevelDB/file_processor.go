
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
}
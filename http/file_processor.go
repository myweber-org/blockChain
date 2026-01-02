package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu          sync.RWMutex
	processed   map[string]bool
	workerCount int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		processed:   make(map[string]bool),
		workerCount: workers,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) ([]string, error) {
	var wg sync.WaitGroup
	results := make(chan string, len(paths))
	errors := make(chan error, len(paths))
	semaphore := make(chan struct{}, fp.workerCount)

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := fp.validatePath(p); err != nil {
				errors <- fmt.Errorf("validation failed for %s: %w", p, err)
				return
			}

			fp.mu.Lock()
			if fp.processed[p] {
				fp.mu.Unlock()
				return
			}
			fp.processed[p] = true
			fp.mu.Unlock()

			content, err := fp.readFile(p)
			if err != nil {
				errors <- err
				return
			}

			processed := fp.transformContent(content)
			results <- processed
		}(path)
	}

	wg.Wait()
	close(results)
	close(errors)

	var resultSlice []string
	for res := range results {
		resultSlice = append(resultSlice, res)
	}

	select {
	case err := <-errors:
		return resultSlice, err
	default:
		return resultSlice, nil
	}
}

func (fp *FileProcessor) validatePath(path string) error {
	if path == "" {
		return errors.New("empty path provided")
	}
	if !filepath.IsAbs(path) {
		return errors.New("path must be absolute")
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access path: %w", err)
	}
	if info.IsDir() {
		return errors.New("path points to directory, file expected")
	}
	return nil
}

func (fp *FileProcessor) readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content []byte
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("read error: %w", err)
		}
		if n == 0 {
			break
		}
		content = append(content, buffer[:n]...)
	}

	return string(content), nil
}

func (fp *FileProcessor) transformContent(content string) string {
	transformed := make([]byte, 0, len(content))
	for i := 0; i < len(content); i++ {
		b := content[i]
		if b >= 'a' && b <= 'z' {
			b = 'a' + (b-'a'+13)%26
		} else if b >= 'A' && b <= 'Z' {
			b = 'A' + (b-'A'+13)%26
		}
		transformed = append(transformed, b)
	}
	return string(transformed)
}

func main() {
	processor := NewFileProcessor(4)
	files := []string{"/tmp/test1.txt", "/tmp/test2.txt"}

	results, err := processor.ProcessFiles(files)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for i, result := range results {
		fmt.Printf("File %d processed: %d characters\n", i+1, len(result))
	}
}
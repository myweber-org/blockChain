package main

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
	if workers < 1 {
		workers = 1
	}
	if batchSize < 1 {
		batchSize = 10
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string, processor func(string) error) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("read directory failed: %w", err)
	}

	fileChan := make(chan string, fp.batchSize)
	errChan := make(chan error, fp.workers)

	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(fileChan, errChan, processor)
	}

	go func() {
		for _, file := range files {
			if !file.IsDir() {
				fullPath := filepath.Join(dirPath, file.Name())
				fileChan <- fullPath
			}
		}
		close(fileChan)
	}()

	go func() {
		fp.wg.Wait()
		close(errChan)
	}()

	var processErrors []error
	for err := range errChan {
		if err != nil {
			processErrors = append(processErrors, err)
		}
	}

	if len(processErrors) > 0 {
		return fmt.Errorf("processing completed with %d errors", len(processErrors))
	}
	return nil
}

func (fp *FileProcessor) worker(files <-chan string, errChan chan<- error, processor func(string) error) {
	defer fp.wg.Done()
	for file := range files {
		if err := processor(file); err != nil {
			errChan <- fmt.Errorf("process %s: %w", file, err)
		}
	}
}

func CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan file: %w", err)
	}
	return lineCount, nil
}

func ValidateFileExtension(filePath string, allowedExts []string) error {
	ext := filepath.Ext(filePath)
	for _, allowed := range allowedExts {
		if ext == allowed {
			return nil
		}
	}
	return errors.New("unsupported file extension")
}

func CreateBackup(originalPath string) (string, error) {
	backupPath := originalPath + ".bak_" + time.Now().Format("20060102_150405")
	
	src, err := os.Open(originalPath)
	if err != nil {
		return "", fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("create backup: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("copy data: %w", err)
	}

	return backupPath, nil
}

func main() {
	processor := NewFileProcessor(4, 20)
	
	dir := "./testfiles"
	err := processor.ProcessDirectory(dir, func(filePath string) error {
		lines, err := CountLines(filePath)
		if err != nil {
			return err
		}
		
		exts := []string{".txt", ".go", ".md"}
		if err := ValidateFileExtension(filePath, exts); err != nil {
			return err
		}
		
		backupPath, err := CreateBackup(filePath)
		if err != nil {
			return err
		}
		
		fmt.Printf("Processed: %s (lines: %d, backup: %s)\n", 
			filepath.Base(filePath), lines, filepath.Base(backupPath))
		return nil
	})
	
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("All files processed successfully")
}
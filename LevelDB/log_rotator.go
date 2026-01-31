
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type LogRotator struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = stat.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	backupPath := fmt.Sprintf("%s.%d", lr.filePath, lr.rotationCount)
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	lr.rotationCount++

	file, err := os.Create(lr.filePath)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 1024*1024) // 1MB max size
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message\n", i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}
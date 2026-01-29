
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
		fileIndex: 0,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := fmt.Sprintf("%s.%d.log", rl.filePath, rl.fileIndex)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentSize < rl.maxSize {
		return nil
	}

	rl.currentFile.Close()
	rl.fileIndex++

	oldFilename := fmt.Sprintf("%s.%d.log", rl.filePath, rl.fileIndex-1)
	newFilename := fmt.Sprintf("%s.%d.log", rl.filePath, rl.fileIndex-1) + "." + time.Now().Format("20060102-150405")
	os.Rename(oldFilename, newFilename)

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err = rl.currentFile.Write(p)
	rl.currentSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) WriteString(s string) (n int, err error) {
	return rl.Write([]byte(s))
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("application", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: This is a sample log message.\n", 
			time.Now().Format(time.RFC3339), i)
		logger.WriteString(message)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
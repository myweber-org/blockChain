package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logger := &RotatingLogger{
		baseName: baseName,
	}

	if err := logger.openNewFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openNewFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

	return nil
}

func (rl *RotatingLogger) compressOldFile() error {
	if rl.currentFile == nil {
		return nil
	}

	oldPath := rl.currentFile.Name()
	newPath := oldPath + ".gz"

	oldFile, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	newFile, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	gzWriter := gzip.NewWriter(newFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return err
	}

	os.Remove(oldPath)
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.compressOldFile(); err != nil {
			return 0, err
		}
		if err := rl.openNewFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	rl.currentSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}
	
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}
	
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}
	
	rotationTime := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", rl.filePath, rotationTime, rl.rotationCount)
	
	err := os.Rename(rl.filePath, backupPath)
	if err != nil {
		return err
	}
	
	rl.rotationCount++
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		rl.mu.Unlock()
		if err := rl.rotate(); err != nil {
			rl.mu.Lock()
			return 0, err
		}
		rl.mu.Lock()
	}
	
	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}
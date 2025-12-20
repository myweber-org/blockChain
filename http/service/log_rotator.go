
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	filePath   string
	currentSize int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: path,
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if _, err := os.Stat(rl.filePath); err == nil {
		os.Rename(rl.filePath, rl.backupPath(0))
	}

	return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.filePath, index)
}

func (rl *RotatingLogger) compressOldLogs() {
	for i := 1; i <= backupCount; i++ {
		path := fmt.Sprintf("%s.%d", rl.filePath, i)
		if _, err := os.Stat(path); err == nil && !strings.HasSuffix(path, ".gz") {
			go func(p string) {
				rl.compressFile(p)
			}(path)
		}
	}
}

func (rl *RotatingLogger) compressFile(path string) {
	// Compression implementation would go here
	// For now just rename to .gz extension
	newPath := path + ".gz"
	os.Rename(path, newPath)
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}

	logger.compressOldLogs()
}
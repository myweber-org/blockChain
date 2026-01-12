package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(logDir string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    fullPath := filepath.Join(logDir, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to stat log file: %w", err)
    }

    return &LogRotator{
        currentFile: file,
        currentSize: info.Size(),
        basePath:    logDir,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return fmt.Errorf("failed to close current log file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    backupPath := filepath.Join(lr.basePath, backupName)
    originalPath := filepath.Join(lr.basePath, logFileName)

    if err := os.Rename(originalPath, backupPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(originalPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    lr.currentFile = file
    lr.currentSize = 0

    go lr.cleanupOldLogs()
    return nil
}

func (lr *LogRotator) cleanupOldLogs() {
    files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i := maxBackupCount; i < len(files); i++ {
        os.Remove(files[i])
    }
}

func (lr *LogRotator) Close() error {
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
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
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}
	
	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return logger, nil
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

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentSize < rl.maxSize {
		return nil
	}
	
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	rl.currentFile.Close()
	
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationCount)
	rl.rotationCount++
	
	if err := os.Rename(rl.filePath, archivePath); err != nil {
		return err
	}
	
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}
	
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	n, err := rl.currentFile.Write(p)
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
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}
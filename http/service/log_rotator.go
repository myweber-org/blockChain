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
	rotationNum int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

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

	rl.currentFile.Close()

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationNum)
	rl.rotationNum++

	if err := os.Rename(rl.filePath, archivePath); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
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

type Rotator struct {
    mu          sync.Mutex
    file        *os.File
    basePath    string
    maxSize     int64
    maxAge      time.Duration
    currentSize int64
    createdAt   time.Time
}

func NewRotator(basePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    dir := filepath.Dir(basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &Rotator{
        file:        file,
        basePath:    basePath,
        maxSize:     maxSize,
        maxAge:      maxAge,
        currentSize: info.Size(),
        createdAt:   time.Now(),
    }, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if err := r.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotateIfNeeded() error {
    needsRotation := false
    var reason string

    if r.maxSize > 0 && r.currentSize >= r.maxSize {
        needsRotation = true
        reason = "size"
    }

    if r.maxAge > 0 && time.Since(r.createdAt) >= r.maxAge {
        needsRotation = true
        reason = "age"
    }

    if !needsRotation {
        return nil
    }

    if err := r.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s_%s", r.basePath, reason, timestamp)

    if err := os.Rename(r.basePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    r.file = file
    r.currentSize = 0
    r.createdAt = time.Now()

    go r.cleanOldFiles()
    return nil
}

func (r *Rotator) cleanOldFiles() {
    if r.maxAge <= 0 {
        return
    }

    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    cutoff := time.Now().Add(-r.maxAge * 2)

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        name := entry.Name()
        if !isBackupFile(name, baseName) {
            continue
        }

        info, err := entry.Info()
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoff) {
            os.Remove(filepath.Join(dir, name))
        }
    }
}

func isBackupFile(name, baseName string) bool {
    return len(name) > len(baseName) && name[:len(baseName)] == baseName
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.file.Close()
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    io.WriteString(rotator, "Test log entry\n")
    fmt.Println("Log rotation system initialized")
}
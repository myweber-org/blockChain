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
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	if info, err := file.Stat(); err == nil {
		rl.currentSize = info.Size()
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	if err := rl.compressOldLogs(); err != nil {
		return err
	}

	if err := rl.cleanupOldBackups(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file) == ".gz" {
			continue
		}

		if err := compressFile(file); err != nil {
			return err
		}
	}
	return nil
}

func compressFile(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	return os.Remove(src)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.gz"))
	if err != nil {
		return err
	}

	if len(files) <= maxBackups {
		return nil
	}

	for i := 0; i < len(files)-maxBackups; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}
	return nil
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

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

type RotatingLogger struct {
    mu            sync.Mutex
    currentFile   *os.File
    basePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile: file,
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
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
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, archivedPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++

    return nil
}

func compressFile(source, target string) error {
    sourceFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    targetFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    gzWriter := gzip.NewWriter(targetFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    return err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
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
	"strconv"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu            sync.Mutex
	basePath      string
	maxSize       int64
	maxBackups    int
	currentSize   int64
	currentFile   *os.File
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxBackups int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:   basePath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	rl.rotationCount++

	if err := rl.cleanupOldBackups(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to cleanup old backups: %v\n", err)
	}

	if err := rl.openCurrentFile(); err != nil {
		return fmt.Errorf("failed to open new log file after rotation: %w", err)
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, baseName+".") && entry.Type().IsRegular() {
			backups = append(backups, name)
		}
	}

	if len(backups) <= rl.maxBackups {
		return nil
	}

	sortBackupsByTimestamp(backups)

	for i := 0; i < len(backups)-rl.maxBackups; i++ {
		path := filepath.Join(dir, backups[i])
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", backups[i], err)
		}
	}

	return nil
}

func sortBackupsByTimestamp(backups []string) {
	for i := 0; i < len(backups); i++ {
		for j := i + 1; j < len(backups); j++ {
			ts1 := extractTimestamp(backups[i])
			ts2 := extractTimestamp(backups[j])
			if ts1 > ts2 {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}
}

func extractTimestamp(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
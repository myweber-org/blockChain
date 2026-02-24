
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (r *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(r.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        err := r.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err != nil {
        return n, err
    }

    r.currentSize += int64(n)
    return n, nil
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

    err := os.Rename(r.basePath, backupPath)
    if err != nil {
        return err
    }

    if r.compressOld {
        compressedPath := backupPath + ".gz"
        err := compressFile(backupPath, compressedPath)
        if err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }

    err = r.cleanOldBackups()
    if err != nil {
        return err
    }

    return r.openCurrentFile()
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (r *LogRotator) cleanOldBackups() error {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backups = append(backups, name)
        }
    }

    if len(backups) <= r.maxBackups {
        return nil
    }

    sortBackups(backups)

    for i := 0; i < len(backups)-r.maxBackups; i++ {
        path := filepath.Join(dir, backups[i])
        os.Remove(path)
    }

    return nil
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return 0
    }

    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }

    timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }

    return timestamp
}

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}
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

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
	maxBackups   int
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

	file, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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
			log.Printf("Failed to rotate log file: %v", err)
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

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, rotatedPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := rl.compressFile(rotatedPath); err != nil {
		log.Printf("Failed to compress log file %s: %v", rotatedPath, err)
	}

	rl.rotationCount++
	if rl.rotationCount > rl.maxBackups {
		rl.cleanupOldBackups()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := rl.basePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Failed to list backup files: %v", err)
		return
	}

	if len(matches) <= rl.maxBackups {
		return
	}

	for i := 0; i < len(matches)-rl.maxBackups; i++ {
		if err := os.Remove(matches[i]); err != nil {
			log.Printf("Failed to remove old backup %s: %v", matches[i], err)
		}
	}
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
		log.Fatalf("Failed to create rotating logger: %v", err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry number %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
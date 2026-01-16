
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
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    currentFile *os.File
}

func NewLogRotator(basePath string, maxSizeMB int, maxFiles int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    err := rotator.initialize()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) initialize() error {
    dir := filepath.Dir(lr.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create log directory: %w", err)
    }

    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return fmt.Errorf("failed to stat log file: %w", err)
    }

    lr.currentFile = file
    lr.currentSize = info.Size()

    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, fmt.Errorf("failed to rotate log: %w", err)
        }
    }

    n, err := lr.currentFile.Write(p)
    if err != nil {
        return n, fmt.Errorf("failed to write to log file: %w", err)
    }

    lr.currentSize += int64(n)
    return n, nil
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return fmt.Errorf("failed to compress rotated log: %w", err)
    }

    file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    lr.currentFile = file
    lr.currentSize = 0

    go lr.cleanupOldFiles()

    return nil
}

func (lr *LogRotator) compressFile(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return fmt.Errorf("failed to open source file: %w", err)
    }
    defer srcFile.Close()

    destPath := srcPath + ".gz"
    destFile, err := os.Create(destPath)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return fmt.Errorf("failed to compress data: %w", err)
    }

    os.Remove(srcPath)

    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    dir := filepath.Dir(lr.basePath)
    baseName := filepath.Base(lr.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var rotatedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            rotatedFiles = append(rotatedFiles, filepath.Join(dir, name))
        }
    }

    if len(rotatedFiles) <= lr.maxFiles {
        return
    }

    filesToDelete := len(rotatedFiles) - lr.maxFiles
    for i := 0; i < filesToDelete; i++ {
        os.Remove(rotatedFiles[i])
    }
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
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: Application event at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation demo completed")
}

package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCount   int
}

func NewRotatingLog(basePath string, maxSize int64) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSize,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) openCurrentFile() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLog) rotate() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile == nil {
        return fmt.Errorf("no current file")
    }

    rl.currentFile.Close()
    rl.fileCount++

    timestamp := time.Now().Format("20060102_150405")
    archiveName := fmt.Sprintf("%s.%s.%d.gz", rl.basePath, timestamp, rl.fileCount)

    oldFile, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    archiveFile, err := os.Create(archiveName)
    if err != nil {
        return err
    }
    defer archiveFile.Close()

    gzWriter := gzip.NewWriter(archiveFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
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

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func (rl *RotatingLog) ListArchives() ([]string, error) {
    pattern := rl.basePath + ".*.gz"
    return filepath.Glob(pattern)
}

func main() {
    logFile, err := NewRotatingLog("app.log", 1024*1024)
    if err != nil {
        fmt.Printf("Failed to create log file: %v\n", err)
        return
    }
    defer logFile.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logFile.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    archives, err := logFile.ListArchives()
    if err != nil {
        fmt.Printf("Error listing archives: %v\n", err)
        return
    }

    fmt.Printf("Created %d archive(s):\n", len(archives))
    for _, archive := range archives {
        fmt.Printf("  %s\n", archive)
    }
}
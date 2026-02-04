
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        sequence: 0,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        lr.file.Close()
    }

    filename := lr.basePath
    if lr.sequence > 0 {
        filename = fmt.Sprintf("%s.%d", lr.basePath, lr.sequence)
    }

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.file = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) rotate() error {
    lr.sequence++
    return lr.openCurrentFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func (lr *LogRotator) CleanOldLogs(maxFiles int) error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    for i := lr.sequence - maxFiles; i >= 0; i-- {
        filename := lr.basePath
        if i > 0 {
            filename = fmt.Sprintf("%s.%d", lr.basePath, i)
        }
        if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
            return err
        }
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message.\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    if err := rotator.CleanOldLogs(5); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }
}
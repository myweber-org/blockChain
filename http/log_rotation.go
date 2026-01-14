package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    file        *os.File
    currentSize int64
    maxSize     int64
    basePath    string
    sequence    int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        maxSize:  maxSize,
        basePath: basePath,
    }
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    path := rl.basePath + ".log"
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.log", rl.basePath, timestamp)
    if err := os.Rename(rl.basePath+".log", rotatedPath); err != nil {
        return err
    }

    go rl.compressOldLog(rotatedPath)

    return rl.openCurrent()
}

func (rl *RotatingLogger) compressOldLog(path string) {
    src, err := os.Open(path)
    if err != nil {
        log.Printf("Failed to open log for compression: %v", err)
        return
    }
    defer src.Close()

    dstPath := path + ".gz"
    dst, err := os.Create(dstPath)
    if err != nil {
        log.Printf("Failed to create compressed file: %v", err)
        return
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        log.Printf("Compression failed: %v", err)
        return
    }

    if err := os.Remove(path); err != nil {
        log.Printf("Failed to remove uncompressed log: %v", err)
    }
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
    logger, err := NewRotatingLogger("app", 1024*1024) // 1MB max size
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    for i := 0; i < 10000; i++ {
        msg := fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())
        if _, err := logger.Write([]byte(msg)); err != nil {
            log.Printf("Write error: %v", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
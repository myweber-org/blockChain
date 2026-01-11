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

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSize,
    }
    err := rl.openCurrentFile()
    return rl, err
}

func (rl *RotatingLog) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
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
    rl.fileCounter = rl.findLatestCounter()
    return nil
}

func (rl *RotatingLog) findLatestCounter() int {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil || len(matches) == 0 {
        return 0
    }

    maxCounter := 0
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        counter, err := strconv.Atoi(parts[len(parts)-2])
        if err == nil && counter > maxCounter {
            maxCounter = counter
        }
    }
    return maxCounter
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
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.fileCounter++
    archivePath := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.fileCounter, time.Now().Format("20060102"))

    if err := rl.compressFile(rl.basePath, archivePath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLog) compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    gzWriter := gzip.NewWriter(destination)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, source)
    return err
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}
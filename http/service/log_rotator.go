
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

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }
    if err := logger.initialize(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (rl *RotatingLogger) initialize() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    rl.fileCounter = rl.findLatestCounter()
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) findLatestCounter() int {
    pattern := rl.basePath + ".*.log"
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

func (rl *RotatingLogger) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileCounter)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
    if rl.currentFile != nil {
        rl.currentFile.Close()
        if err := rl.compressOldFile(); err != nil {
            return err
        }
    }
    rl.fileCounter++
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldFile() error {
    oldName := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileCounter-1)
    newName := oldName + ".gz"
    oldFile, err := os.Open(oldName)
    if err != nil {
        return err
    }
    defer oldFile.Close()
    newFile, err := os.Create(newName)
    if err != nil {
        return err
    }
    defer newFile.Close()
    gzWriter := gzip.NewWriter(newFile)
    defer gzWriter.Close()
    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }
    return os.Remove(oldName)
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
    logger, err := NewRotatingLogger("/var/log/myapp/application", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}
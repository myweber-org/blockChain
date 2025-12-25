package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
    logDir      = "logs"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    basePath := filepath.Join(logDir, name)
    rl := &RotatingLogger{baseName: basePath}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    path := rl.baseName + ".log"
    file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
    timestamp := time.Now().Format("20060102-150405")
    newPath := fmt.Sprintf("%s-%s.log", rl.baseName, timestamp)

    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    if err := os.Rename(rl.baseName+".log", newPath); err != nil {
        return err
    }

    if err := rl.openCurrent(); err != nil {
        return err
    }

    rl.cleanupOld()
    return nil
}

func (rl *RotatingLogger) cleanupOld() {
    pattern := rl.baseName + "-*.log"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    toDelete := matches[:len(matches)-maxBackups]
    for _, path := range toDelete {
        os.Remove(path)
    }
}

func (rl *RotatingLogger) Close() error {
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

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    FilePath    string
    MaxSize     int64
    MaxFiles    int
    RotateEvery time.Duration
    lastRotate  time.Time
}

func NewLogRotator(filePath string, maxSize int64, maxFiles int, rotateEvery time.Duration) *LogRotator {
    return &LogRotator{
        FilePath:    filePath,
        MaxSize:     maxSize,
        MaxFiles:    maxFiles,
        RotateEvery: rotateEvery,
        lastRotate:  time.Now(),
    }
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if err := lr.checkRotation(); err != nil {
        return 0, err
    }

    file, err := os.OpenFile(lr.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    return file.Write(p)
}

func (lr *LogRotator) checkRotation() error {
    now := time.Now()
    shouldRotate := false

    if lr.RotateEvery > 0 && now.Sub(lr.lastRotate) >= lr.RotateEvery {
        shouldRotate = true
        lr.lastRotate = now
    }

    if !shouldRotate && lr.MaxSize > 0 {
        if info, err := os.Stat(lr.FilePath); err == nil && info.Size() >= lr.MaxSize {
            shouldRotate = true
        }
    }

    if shouldRotate {
        return lr.performRotation()
    }
    return nil
}

func (lr *LogRotator) performRotation() error {
    for i := lr.MaxFiles - 1; i > 0; i-- {
        oldName := fmt.Sprintf("%s.%d", lr.FilePath, i)
        newName := fmt.Sprintf("%s.%d", lr.FilePath, i+1)

        if _, err := os.Stat(oldName); err == nil {
            if err := os.Rename(oldName, newName); err != nil {
                return err
            }
        }
    }

    if _, err := os.Stat(lr.FilePath); err == nil {
        backupName := fmt.Sprintf("%s.1", lr.FilePath)
        return os.Rename(lr.FilePath, backupName)
    }

    return nil
}

func (lr *LogRotator) Cleanup() error {
    for i := lr.MaxFiles + 1; ; i++ {
        fileName := fmt.Sprintf("%s.%d", lr.FilePath, i)
        if _, err := os.Stat(fileName); os.IsNotExist(err) {
            break
        }
        if err := os.Remove(fileName); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    rotator := NewLogRotator(
        "/var/log/app.log",
        10*1024*1024,
        5,
        time.Hour*24,
    )

    message := fmt.Sprintf("[%s] Application started\n", time.Now().Format(time.RFC3339))
    _, err := rotator.Write([]byte(message))
    if err != nil {
        fmt.Printf("Write error: %v\n", err)
        return
    }

    if err := rotator.Cleanup(); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }

    fmt.Println("Log rotation completed successfully")
}
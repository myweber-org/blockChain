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
}package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    baseName    string
    sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
        sequence: 0,
    }

    if err := rl.openNewFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    filename := filepath.Join(logDir, fmt.Sprintf("%s.log", rl.baseName))
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) rotateIfNeeded() error {
    if rl.currentSize < maxFileSize {
        return nil
    }

    oldPath := filepath.Join(logDir, fmt.Sprintf("%s.log", rl.baseName))
    newPath := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, rl.sequence))

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    rl.sequence++
    if rl.sequence > maxBackups {
        rl.cleanupOldFiles()
    }

    return rl.openNewFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
    for i := 0; i <= rl.sequence-maxBackups; i++ {
        oldFile := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, i))
        os.Remove(oldFile)
    }
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if err := rl.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err = rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
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
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(10 * time.Millisecond)
    }
}

package main

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
)

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    baseName    string
    fileSize    int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, err
    }

    dir := filepath.Dir(absPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        filePath: absPath,
        baseName: filepath.Base(absPath),
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.fileSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.fileSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }

    if err := rl.openCurrentFile(); err != nil {
        return err
    }

    rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    dir := filepath.Dir(rl.filePath)
    pattern := rl.baseName + ".*"

    matches, err := filepath.Glob(filepath.Join(dir, pattern))
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("logs/app.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: This is a test log message for rotation testing", i)
        time.Sleep(10 * time.Millisecond)
    }
}
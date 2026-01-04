
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    sequence    int
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: basePath,
        sequence: 0,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
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
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    oldPath := rl.currentFile.Name()
    compressedPath := fmt.Sprintf("%s.%d.gz", oldPath, time.Now().Unix())

    if err := compressFile(oldPath, compressedPath); err != nil {
        return err
    }

    os.Remove(oldPath)

    rl.sequence++
    if rl.sequence > maxBackups {
        rl.cleanOldBackups()
    }

    return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (rl *RotatingLogger) openCurrentFile() error {
    path := fmt.Sprintf("%s.log", rl.basePath)
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = stat.Size()
    return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := fmt.Sprintf("%s.log.*.gz", rl.basePath)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        for i := 0; i < len(matches)-maxBackups; i++ {
            os.Remove(matches[i])
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("application")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
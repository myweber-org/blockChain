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
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type RotatingLogger struct {
    filePath string
    current  *os.File
    size     int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        filePath: path,
        current:  file,
        size:     info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.size+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    rl.size += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.current.Close()

    for i := maxBackupFiles - 1; i >= 0; i-- {
        oldPath := rl.backupPath(i)
        newPath := rl.backupPath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            if i == maxBackupFiles-1 {
                os.Remove(oldPath)
            } else {
                os.Rename(oldPath, newPath)
            }
        }
    }

    compressedPath := rl.backupPath(0) + ".gz"
    if err := compressFile(rl.filePath, compressedPath); err != nil {
        return err
    }

    os.Remove(rl.filePath)

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.current = file
    rl.size = 0
    return nil
}

func (rl *RotatingLogger) backupPath(index int) string {
    if index == 0 {
        return rl.filePath + "." + time.Now().Format("20060102150405")
    }
    return fmt.Sprintf("%s.%d", rl.filePath, index)
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    target, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer target.Close()

    gz := gzip.NewWriter(target)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
}

func (rl *RotatingLogger) Close() error {
    return rl.current.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed. Check app.log* files in", filepath.Dir("app.log"))
}
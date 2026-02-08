package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    file        *os.File
}

func NewRotatingLogger(filePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filePath: filePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rl.openFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openFile() error {
    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) rotate() error {
    rl.file.Close()

    for i := rl.maxFiles - 1; i > 0; i-- {
        oldPath := rl.backupPath(i)
        newPath := rl.backupPath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    if err := os.Rename(rl.filePath, rl.backupPath(1)); err != nil {
        return err
    }

    return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
    if index == 0 {
        return rl.filePath
    }
    return rl.filePath + "." + strconv.Itoa(index)
}

func (rl *RotatingLogger) cleanup() error {
    dir := filepath.Dir(rl.filePath)
    base := filepath.Base(rl.filePath)

    files, err := filepath.Glob(filepath.Join(dir, base+".*"))
    if err != nil {
        return err
    }

    var backupFiles []string
    for _, file := range files {
        parts := strings.Split(file, ".")
        if len(parts) > 0 {
            if _, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
                backupFiles = append(backupFiles, file)
            }
        }
    }

    sort.Strings(backupFiles)

    if len(backupFiles) > rl.maxFiles {
        for _, file := range backupFiles[rl.maxFiles:] {
            os.Remove(file)
        }
    }

    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
        if err := rl.cleanup(); err != nil {
            fmt.Printf("Cleanup error: %v\n", err)
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
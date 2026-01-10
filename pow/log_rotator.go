
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
    mu            sync.Mutex
    basePath      string
    currentSize   int64
    maxSize       int64
    currentFile   *os.File
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int64) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSizeMB * 1024 * 1024,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
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
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    rl.rotationCount++
    if err := rl.compressOldLog(rotatedPath); err != nil {
        fmt.Printf("Compression failed: %v\n", err)
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldLog(path string) error {
    if !strings.HasSuffix(path, ".gz") {
        srcFile, err := os.Open(path)
        if err != nil {
            return err
        }
        defer srcFile.Close()

        destFile, err := os.Create(path + ".gz")
        if err != nil {
            return err
        }
        defer destFile.Close()

        gzWriter := gzip.NewWriter(destFile)
        defer gzWriter.Close()

        if _, err := io.Copy(gzWriter, srcFile); err != nil {
            return err
        }

        if err := os.Remove(path); err != nil {
            return err
        }
    }
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles(maxFiles int) error {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var logFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName) && (strings.HasSuffix(name, ".gz") || strings.Contains(name, ".")) {
            logFiles = append(logFiles, filepath.Join(dir, name))
        }
    }

    if len(logFiles) <= maxFiles {
        return nil
    }

    for i := 0; i < len(logFiles)-maxFiles; i++ {
        if err := os.Remove(logFiles[i]); err != nil {
            return err
        }
    }
    return nil
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
    logger, err := NewRotatingLogger("/var/log/app/app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Iteration %d: Processing data chunk %d\n",
            time.Now().Format(time.RFC3339),
            i,
            i*256)
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }

        if i%100 == 0 {
            if err := logger.cleanupOldFiles(5); err != nil {
                fmt.Printf("Cleanup error: %v\n", err)
            }
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed. Check /var/log/app/ directory.")
}
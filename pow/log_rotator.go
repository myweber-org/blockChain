
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

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
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
    if rl.currentFile != nil {
        rl.currentFile.Close()
        
        // Compress the rotated file
        if err := rl.compressFile(rl.currentFile.Name()); err != nil {
            return err
        }
        
        // Clean old backups
        rl.cleanOldBackups()
    }
    
    rl.sequence++
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.sequence)
    f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    
    rl.currentFile = f
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) compressFile(src string) error {
    dst := src + ".gz"
    
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()
    
    gz := gzip.NewWriter(dstFile)
    defer gz.Close()
    
    _, err = io.Copy(gz, srcFile)
    if err != nil {
        return err
    }
    
    // Remove original after compression
    return os.Remove(src)
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := rl.basePath + ".*.log.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) > maxBackups {
        // Sort by modification time (oldest first)
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
    logger, err := NewRotatingLogger("app")
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    // Simulate log writing
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Iteration %d: Processing data chunk\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        
        // Simulate some delay
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation completed")
}
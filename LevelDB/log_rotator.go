
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

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compress bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compress,
    }
    
    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.currentSize+int64(len(p)) > r.maxSize {
        err := r.rotate()
        if err != nil {
            return 0, err
        }
    }
    
    n, err := r.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    
    r.currentSize += int64(n)
    return n, nil
}

func (r *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    
    err := os.Rename(r.basePath, backupPath)
    if err != nil {
        return err
    }
    
    if r.compressOld {
        compressedPath := backupPath + ".gz"
        err := compressFile(backupPath, compressedPath)
        if err != nil {
            return err
        }
        os.Remove(backupPath)
        backupPath = compressedPath
    }
    
    err = r.openCurrentFile()
    if err != nil {
        return err
    }
    
    r.cleanupOldBackups()
    return nil
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()
    
    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()
    
    gz := gzip.NewWriter(destination)
    defer gz.Close()
    
    _, err = io.Copy(gz, source)
    return err
}

func (r *LogRotator) cleanupOldBackups() {
    if r.maxBackups <= 0 {
        return
    }
    
    pattern := r.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    var backupFiles []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") || isTimestampBackup(match, r.basePath) {
            backupFiles = append(backupFiles, match)
        }
    }
    
    if len(backupFiles) <= r.maxBackups {
        return
    }
    
    sortBackupsByTime(backupFiles)
    
    for i := 0; i < len(backupFiles)-r.maxBackups; i++ {
        os.Remove(backupFiles[i])
    }
}

func isTimestampBackup(path, basePath string) bool {
    suffix := strings.TrimPrefix(path, basePath+".")
    if len(suffix) != 14 {
        return false
    }
    
    _, err := strconv.ParseInt(suffix, 10, 64)
    return err == nil
}

func sortBackupsByTime(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            timeI := extractTimestamp(files[i])
            timeJ := extractTimestamp(files[j])
            if timeI.After(timeJ) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(path string) time.Time {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    
    for _, part := range parts {
        if len(part) == 14 {
            t, err := time.Parse("20060102150405", part)
            if err == nil {
                return t
            }
        }
    }
    
    info, err := os.Stat(path)
    if err == nil {
        return info.ModTime()
    }
    
    return time.Time{}
}

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}
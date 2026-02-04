
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
    if err == nil {
        r.currentSize += int64(n)
    }
    
    return n, err
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102150405")
    backupPath := r.basePath + "." + timestamp
    
    err := os.Rename(r.basePath, backupPath)
    if err != nil {
        return err
    }
    
    err = r.openCurrentFile()
    if err != nil {
        return err
    }
    
    go r.manageBackups(backupPath)
    
    return nil
}

func (r *LogRotator) manageBackups(backupPath string) {
    if r.compressOld {
        compressedPath := backupPath + ".gz"
        err := compressFile(backupPath, compressedPath)
        if err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }
    
    r.cleanupOldBackups()
}

func compressFile(src, dst string) error {
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
    
    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()
    
    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (r *LogRotator) cleanupOldBackups() {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)
    
    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }
    
    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backups = append(backups, name)
        }
    }
    
    if len(backups) <= r.maxBackups {
        return
    }
    
    sortBackups(backups)
    
    for i := r.maxBackups; i < len(backups); i++ {
        os.Remove(filepath.Join(dir, backups[i]))
    }
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return 0
    }
    
    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }
    
    timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    
    return timestamp
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
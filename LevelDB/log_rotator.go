
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
    mu           sync.Mutex
    basePath     string
    currentFile  *os.File
    maxSize      int64
    currentSize  int64
    maxBackups   int
    compress     bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxBackups int, compress bool) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        basePath:    basePath,
        currentFile: file,
        maxSize:     maxSize,
        currentSize: info.Size(),
        maxBackups:  maxBackups,
        compress:    compress,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize && rl.maxSize > 0 {
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

    timestamp := time.Now().Format("2006-01-02_15-04-05")
    backupPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0

    go rl.manageBackups(backupPath)

    return nil
}

func (rl *RotatingLogger) manageBackups(backupPath string) {
    dir := filepath.Dir(backupPath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && name != baseName {
            backups = append(backups, name)
        }
    }

    if len(backups) > rl.maxBackups {
        sortBackups(backups)
        toRemove := backups[:len(backups)-rl.maxBackups]

        for _, backup := range toRemove {
            path := filepath.Join(dir, backup)
            if rl.compress && !strings.HasSuffix(path, ".gz") {
                if err := compressFile(path); err == nil {
                    os.Remove(path)
                }
            } else {
                os.Remove(path)
            }
        }
    }
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) < extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(name string) int64 {
    parts := strings.Split(name, ".")
    if len(parts) < 2 {
        return 0
    }

    timestampStr := parts[len(parts)-1]
    timestampStr = strings.ReplaceAll(timestampStr, "-", "")
    timestampStr = strings.ReplaceAll(timestampStr, "_", "")

    ts, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return ts
}

func compressFile(src string) error {
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

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5, true)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation example completed")
}

package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    baseName    string
    sequence    int
}

func NewRotatingWriter(baseName string) (*RotatingWriter, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    w := &RotatingWriter{
        baseName: baseName,
        sequence: 0,
    }

    if err := w.openCurrentFile(); err != nil {
        return nil, err
    }

    return w, nil
}

func (w *RotatingWriter) openCurrentFile() error {
    filename := filepath.Join(logDir, fmt.Sprintf("%s.log", w.baseName))
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.currentFile = file
    w.currentSize = stat.Size()
    return nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.currentFile.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    w.sequence++
    if w.sequence > maxBackups {
        w.sequence = 1
    }

    oldName := filepath.Join(logDir, fmt.Sprintf("%s.log", w.baseName))
    newName := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", w.baseName, w.sequence))

    if err := os.Rename(oldName, newName); err != nil && !os.IsNotExist(err) {
        return err
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        os.Exit(1)
    }
    defer writer.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := writer.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    fileCount   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
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
        currentFile: file,
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        fileCount:   0,
    }, nil
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
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, archivedPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.fileCount++

    if rl.fileCount > 10 {
        rl.cleanOldArchives()
    }

    return nil
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLogger) cleanOldArchives() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > 10 {
        filesToDelete := matches[:len(matches)-10]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
}
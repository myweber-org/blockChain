package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    fileIndex   int
}

func NewRotatingWriter(basePath string) (*RotatingWriter, error) {
    w := &RotatingWriter{
        basePath: basePath,
    }
    if err := w.openCurrentFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) openCurrentFile() error {
    path := w.basePath
    if w.fileIndex > 0 {
        path = fmt.Sprintf("%s.%d", w.basePath, w.fileIndex)
    }

    file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (w *RotatingWriter) rotate() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    w.fileIndex++
    if w.fileIndex > maxBackups {
        w.fileIndex = 0
        for i := 1; i <= maxBackups; i++ {
            oldPath := fmt.Sprintf("%s.%d", w.basePath, i)
            os.Remove(oldPath)
        }
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.currentFile.Write(p)
    if err != nil {
        return n, err
    }

    w.currentSize += int64(n)
    return n, nil
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
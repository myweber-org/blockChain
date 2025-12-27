
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
)

type RotatingWriter struct {
    mu          sync.Mutex
    current     *os.File
    maxSize     int64
    basePath    string
    currentSize int64
    fileIndex   int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
    w := &RotatingWriter{
        maxSize:  maxSize,
        basePath: basePath,
    }
    if err := w.rotate(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.currentSize+int64(len(p)) > w.maxSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.current.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if w.current != nil {
        w.current.Close()
    }

    w.fileIndex++
    filename := fmt.Sprintf("%s.%d", w.basePath, w.fileIndex)
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }

    w.current = file
    w.currentSize = 0
    return nil
}

func (w *RotatingWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    if w.current != nil {
        return w.current.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log", 1024*1024)
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()

    for i := 0; i < 10000; i++ {
        line := fmt.Sprintf("Log entry %d: Some sample log data here\n", i)
        writer.Write([]byte(line))
    }
    fmt.Println("Log rotation test completed")
}

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
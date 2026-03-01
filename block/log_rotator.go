package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    mu          sync.Mutex
    file        *os.File
    basePath    string
    maxSize     int64
    rotateEvery time.Duration
    lastRotate  time.Time
    currentSize int64
}

func NewRotator(basePath string, maxSize int64, rotateEvery time.Duration) (*Rotator, error) {
    r := &Rotator{
        basePath:    basePath,
        maxSize:     maxSize,
        rotateEvery: rotateEvery,
        lastRotate:  time.Now(),
    }
    if err := r.openFile(); err != nil {
        return nil, err
    }
    return r, nil
}

func (r *Rotator) openFile() error {
    dir := filepath.Dir(r.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    f, err := os.OpenFile(r.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.file = f
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) rotate() error {
    timestamp := time.Now().Format("20060102_150405")
    newPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    if err := os.Rename(r.basePath, newPath); err != nil {
        return err
    }
    r.file.Close()
    if err := r.openFile(); err != nil {
        return err
    }
    r.lastRotate = time.Now()
    r.currentSize = 0
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    shouldRotate := false

    if r.maxSize > 0 && r.currentSize+int64(len(p)) > r.maxSize {
        shouldRotate = true
    }
    if r.rotateEvery > 0 && now.Sub(r.lastRotate) >= r.rotateEvery {
        shouldRotate = true
    }

    if shouldRotate {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("/var/log/myapp/app.log", 10*1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
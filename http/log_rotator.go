
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
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "application.log"
)

type LogRotator struct {
    currentSize int64
    file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    lr.file.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)

    if err := os.Rename(logFileName, backupName); err != nil {
        return err
    }

    if err := lr.compressBackup(backupName); err != nil {
        return err
    }

    file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0

    go lr.cleanupOldBackups()

    return nil
}

func (lr *LogRotator) compressBackup(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dest, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dest.Close()

    gz := gzip.NewWriter(dest)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(filename)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := logFileName + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupCount {
        return
    }

    for i := 0; i < len(matches)-maxBackupCount; i++ {
        os.Remove(matches[i])
    }
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
        rotator.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
}package main

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

	filename := fmt.Sprintf("%s.%d", w.basePath, w.fileIndex)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	w.current = file
	w.currentSize = 0
	w.fileIndex++
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
	writer, err := NewRotatingWriter("app.log", 1024*10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		line := fmt.Sprintf("Log entry number %d: Some sample log data here.\n", i)
		if _, err := writer.Write([]byte(line)); err != nil {
			fmt.Fprintf(os.Stderr, "Write failed: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed. Check app.log.* files.")
}
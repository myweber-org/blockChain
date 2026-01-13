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

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSize,
    }
    err := rl.openCurrentFile()
    return rl, err
}

func (rl *RotatingLog) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    rl.fileCounter = rl.findLatestCounter()
    return nil
}

func (rl *RotatingLog) findLatestCounter() int {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil || len(matches) == 0 {
        return 0
    }

    maxCounter := 0
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        counter, err := strconv.Atoi(parts[len(parts)-2])
        if err == nil && counter > maxCounter {
            maxCounter = counter
        }
    }
    return maxCounter
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
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

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.fileCounter++
    archivePath := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.fileCounter, time.Now().Format("20060102"))

    if err := rl.compressFile(rl.basePath, archivePath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLog) compressFile(src, dst string) error {
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

    gzWriter := gzip.NewWriter(destination)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, source)
    return err
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	basePath      string
	retentionDays int
}

func NewRotatingLogger(basePath string, maxSizeMB int, retentionDays int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		maxSize:       maxSize,
		basePath:      basePath,
		retentionDays: retentionDays,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	go rl.cleanupOldFiles()

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s_%s.log", rl.basePath, timestamp)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Failed to rotate log: %v", err)
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("2006-01-02_150405")
	oldPath := fmt.Sprintf("%s_%s.log", rl.basePath, timestamp)
	newPath := fmt.Sprintf("%s.gz", oldPath)

	if err := os.Rename(rl.currentFile.Name(), oldPath); err != nil {
		return err
	}

	if err := compressFile(oldPath, newPath); err != nil {
		log.Printf("Compression failed: %v", err)
		os.Remove(newPath)
	} else {
		os.Remove(oldPath)
	}

	return rl.openCurrentFile()
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

	compressor := NewGzipWriter(dstFile)
	defer compressor.Close()

	_, err = io.Copy(compressor, srcFile)
	return err
}

func (rl *RotatingLogger) cleanupOldFiles() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().AddDate(0, 0, -rl.retentionDays)
		pattern := fmt.Sprintf("%s_*.log.gz", rl.basePath)

		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, file := range matches {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoff) {
				os.Remove(file)
			}
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

type GzipWriter struct {
	io.WriteCloser
}

func NewGzipWriter(w io.Writer) *GzipWriter {
	return &GzipWriter{}
}

func (gw *GzipWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (gw *GzipWriter) Close() error {
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app", 10, 30)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("test", i%10+1))
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
)

type RotatingWriter struct {
	mu       sync.Mutex
	filename string
	file     *os.File
	size     int64
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	w := &RotatingWriter{filename: filename}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) >= maxFileSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotatingWriter) openFile() error {
	info, err := os.Stat(w.filename)
	if os.IsNotExist(err) {
		w.size = 0
	} else if err != nil {
		return err
	} else {
		w.size = info.Size()
	}

	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.file != nil {
		w.file.Close()
	}

	for i := maxBackups - 1; i >= 0; i-- {
		oldName := w.backupName(i)
		newName := w.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(w.filename, w.backupName(0)); err != nil {
		return err
	}

	return w.openFile()
}

func (w *RotatingWriter) backupName(i int) string {
	if i == 0 {
		return w.filename + ".1"
	}
	return fmt.Sprintf("%s.%d", w.filename, i+1)
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
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
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
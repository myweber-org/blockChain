
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
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDay := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDay {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLog(); err != nil {
				log.Printf("Failed to compress old log: %v", err)
			}
			if err := rl.cleanupOldBackups(); err != nil {
				log.Printf("Failed to cleanup old backups: %v", err)
			}
		}

		newPath := rl.getLogPath(now)
		if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			return fmt.Errorf("create log directory: %w", err)
		}

		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("open log file: %w", err)
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			return fmt.Errorf("stat log file: %w", err)
		}

		rl.file = file
		rl.size = info.Size()
		rl.currentDay = currentDay
	}
	return nil
}

func (rl *RotatingLogger) getLogPath(t time.Time) string {
	if rl.file != nil && rl.size < maxFileSize && rl.currentDay == t.Format("2006-01-02") {
		return rl.file.Name()
	}

	base := filepath.Base(rl.basePath)
	ext := filepath.Ext(rl.basePath)
	name := base[:len(base)-len(ext)]

	for i := 0; i < 1000; i++ {
		var suffix string
		if i > 0 {
			suffix = fmt.Sprintf(".%d", i)
		}
		path := filepath.Join(
			filepath.Dir(rl.basePath),
			t.Format("2006-01-02"),
			fmt.Sprintf("%s%s%s", name, suffix, ext),
		)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}
	}
	return rl.basePath
}

func (rl *RotatingLogger) compressOldLog() error {
	if rl.file == nil {
		return nil
	}

	oldPath := rl.file.Name()
	compressedPath := oldPath + ".gz"

	src, err := os.Open(oldPath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(compressedPath)
	if err != nil {
		return fmt.Errorf("create compressed file: %w", err)
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return fmt.Errorf("compress data: %w", err)
	}

	if err := os.Remove(oldPath); err != nil {
		return fmt.Errorf("remove old file: %w", err)
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	dir := filepath.Dir(rl.basePath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read log directory: %w", err)
	}

	var backupDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			backupDirs = append(backupDirs, entry.Name())
		}
	}

	if len(backupDirs) > backupCount {
		for i := 0; i < len(backupDirs)-backupCount; i++ {
			oldDir := filepath.Join(dir, backupDirs[i])
			if err := os.RemoveAll(oldDir); err != nil {
				return fmt.Errorf("remove old backup %s: %w", oldDir, err)
			}
		}
	}

	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs/app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
		time.Sleep(100 * time.Millisecond)
	}
}
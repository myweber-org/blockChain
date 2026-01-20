
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	if info, err := file.Stat(); err == nil {
		rl.currentSize = info.Size()
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return err
	}

	if len(files) >= maxBackups {
		oldest := files[0]
		for _, f := range files[1:] {
			if fi, _ := os.Stat(f); fi != nil {
				if ofi, _ := os.Stat(oldest); ofi != nil && fi.ModTime().Before(ofi.ModTime()) {
					oldest = f
				}
			}
		}
		if err := rl.compressAndRemove(oldest); err != nil {
			return err
		}
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressAndRemove(filename string) error {
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

	return os.Remove(filename)
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	maxBackups  int
	currentSize int64
	file        *os.File
}

func NewLogRotator(filePath string, maxSize int64, maxBackups int) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath:   filePath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	if err := rotator.openFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openFile() error {
	dir := filepath.Dir(lr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to get file info: %w", err)
	}

	lr.file = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	n, err := lr.file.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := lr.compressFile(backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to compress %s: %v\n", backupPath, err)
	}

	if err := lr.cleanupOldBackups(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to cleanup old backups: %v\n", err)
	}

	if err := lr.openFile(); err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	return nil
}

func (lr *LogRotator) compressFile(sourcePath string) error {
	destPath := sourcePath + ".gz"

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	compressor := NewGzipWriter(destFile)
	defer compressor.Close()

	if _, err := io.Copy(compressor, sourceFile); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to compress data: %w", err)
	}

	if err := os.Remove(sourcePath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to remove uncompressed file %s: %v\n", sourcePath, err)
	}

	return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
	dir := filepath.Dir(lr.filePath)
	baseName := filepath.Base(lr.filePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, baseName+".") && (strings.HasSuffix(name, ".gz") || !strings.Contains(name, ".gz")) {
			backups = append(backups, name)
		}
	}

	if len(backups) <= lr.maxBackups {
		return nil
	}

	backupsToRemove := backups[:len(backups)-lr.maxBackups]
	for _, backup := range backupsToRemove {
		path := filepath.Join(dir, backup)
		if err := os.Remove(path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove old backup %s: %v\n", path, err)
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 1024*1024, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation demonstration completed")
}package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewRotatingWriter(baseName string) (*RotatingWriter, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    basePath := filepath.Join(logDir, baseName)
    w := &RotatingWriter{basePath: basePath}
    if err := w.openCurrent(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) openCurrent() error {
    path := w.basePath + ".log"
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.currentFile = file
    w.currentSize = info.Size()
    return nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.currentFile.Write(p)
    w.currentSize += int64(n)
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    newPath := fmt.Sprintf("%s_%s.log", w.basePath, timestamp)
    oldPath := w.basePath + ".log"

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    if err := w.cleanupOld(); err != nil {
        log.Printf("Cleanup error: %v", err)
    }

    return w.openCurrent()
}

func (w *RotatingWriter) cleanupOld() error {
    pattern := w.basePath + "_*.log"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, path := range toDelete {
            if err := os.Remove(path); err != nil {
                return err
            }
        }
    }
    return nil
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
        log.Fatal(err)
    }
    defer writer.Close()

    logger := log.New(io.MultiWriter(os.Stdout, writer), "", log.LstdFlags)

    for i := 0; i < 1000; i++ {
        logger.Printf("Log entry %d: Application event occurred", i)
        time.Sleep(10 * time.Millisecond)
    }
}
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
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
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
	return nil
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	rl.rotationCount++
	archiveName := fmt.Sprintf("%s.%d-%s.gz", 
		rl.basePath, 
		rl.rotationCount, 
		time.Now().Format("20060102-150405"))

	if err := rl.compressFile(rl.basePath, archiveName); err != nil {
		return err
	}

	if err := os.Truncate(rl.basePath, 0); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/app/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application event recorded\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogRotator struct {
	FilePath    string
	MaxSize     int64
	BackupCount int
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) *LogRotator {
	return &LogRotator{
		FilePath:    filePath,
		MaxSize:     maxSize,
		BackupCount: backupCount,
	}
}

func (lr *LogRotator) RotateIfNeeded() error {
	info, err := os.Stat(lr.FilePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if info.Size() < lr.MaxSize {
		return nil
	}

	dir := filepath.Dir(lr.FilePath)
	baseName := filepath.Base(lr.FilePath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	for i := lr.BackupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		if i == 0 {
			oldPath = lr.FilePath
		} else {
			oldPath = filepath.Join(dir, fmt.Sprintf("%s.%d%s", nameWithoutExt, i, ext))
		}
		newPath = filepath.Join(dir, fmt.Sprintf("%s.%d%s", nameWithoutExt, i+1, ext))

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
			}
		}
	}

	file, err := os.Create(lr.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	file.Close()

	log.Printf("Log rotated: %s", lr.FilePath)
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if err := lr.RotateIfNeeded(); err != nil {
		return 0, err
	}

	file, err := os.OpenFile(lr.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	return file.Write(p)
}

func main() {
	rotator := NewLogRotator("/var/log/myapp/app.log", 1024*1024, 5)

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			log.Fatalf("Failed to write log: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
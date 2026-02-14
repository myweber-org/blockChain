package logutil

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Rotator struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationTime time.Time
}

func NewRotator(basePath string, maxSizeMB int) (*Rotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &Rotator{
		file:         file,
		basePath:     basePath,
		maxSize:      maxSize,
		currentSize:  info.Size(),
		rotationTime: time.Now(),
	}, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentSize+int64(len(p)) > r.maxSize || time.Since(r.rotationTime).Hours() >= 24 {
		if err := r.rotate(); err != nil {
			return 0, fmt.Errorf("rotation failed: %w", err)
		}
	}

	n, err := r.file.Write(p)
	if err == nil {
		r.currentSize += int64(n)
	}
	return n, err
}

func (r *Rotator) rotate() error {
	if err := r.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", r.basePath, timestamp)

	if err := compressFile(r.basePath, archivePath); err != nil {
		return fmt.Errorf("compression failed: %w", err)
	}

	if err := os.Remove(r.basePath); err != nil {
		return fmt.Errorf("failed to remove old log file: %w", err)
	}

	file, err := os.OpenFile(r.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	r.file = file
	r.currentSize = 0
	r.rotationTime = time.Now()
	return nil
}

func compressFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	_, err = io.Copy(gz, source)
	return err
}

func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.file.Close()
}

func (r *Rotator) Flush() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.file.Sync()
}

func (r *Rotator) CurrentSize() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.currentSize
}

func (r *Rotator) RotationTime() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rotationTime
}

func (r *Rotator) FilePath() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.basePath
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	rl := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}
	
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}
	
	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationCount)
	
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}
	
	rl.rotationCount++
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log: %w", err)
		}
	}
	
	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, fmt.Errorf("failed to write to log: %w", err)
	}
	
	rl.currentSize += int64(n)
	return n, nil
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
	logger, err := NewRotatingLogger("./logs/app.log", 10)
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
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}
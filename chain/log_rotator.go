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
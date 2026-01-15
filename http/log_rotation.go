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

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	basePath   string
	maxSize    int64
	currentSize int64
	sequence   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		sequence: 0,
	}
	if err := rl.openOrCreate(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openOrCreate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	path := rl.basePath
	if rl.sequence > 0 {
		path = fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
		rl.file = nil
	}

	oldPath := rl.basePath
	if rl.sequence > 0 {
		oldPath = fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
	}

	rl.sequence++
	if err := rl.openOrCreate(); err != nil {
		return err
	}

	go rl.compressOldLog(oldPath)
	return nil
}

func (rl *RotatingLogger) compressOldLog(path string) {
	src, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open log for compression: %v", err)
		return
	}
	defer src.Close()

	destPath := path + ".gz"
	dest, err := os.Create(destPath)
	if err != nil {
		log.Printf("Failed to create compressed file: %v", err)
		return
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		log.Printf("Compression failed: %v", err)
		os.Remove(destPath)
		return
	}

	if err := os.Remove(path); err != nil {
		log.Printf("Failed to remove original log: %v", err)
	}
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
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed. Check app.log* files.")
}
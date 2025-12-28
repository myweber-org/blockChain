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
	maxLogSize   = 10 * 1024 * 1024 // 10MB
	logDir       = "./logs"
	currentLog   = "app.log"
	timeFormat   = "2006-01-02_15-04-05"
)

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	size     int64
	basePath string
}

func NewRotatingLogger() (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	path := filepath.Join(logDir, currentLog)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		file:     file,
		size:     info.Size(),
		basePath: path,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.size += int64(n)
	if rl.size >= maxLogSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}

	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format(timeFormat)
	archiveName := fmt.Sprintf("app_%s.log.gz", timestamp)
	archivePath := filepath.Join(logDir, archiveName)

	oldLog, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer oldLog.Close()

	archive, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	gzWriter := gzip.NewWriter(archive)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldLog); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	newFile, err := os.Create(rl.basePath)
	if err != nil {
		return err
	}

	rl.file = newFile
	rl.size = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Simulating application activity", i)
		time.Sleep(10 * time.Millisecond)
	}
}
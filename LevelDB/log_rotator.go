package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logDir         = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
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

	timestamp := time.Now().Format("20060102-150405")
	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		log.Printf("Cleanup failed: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := filepath.Join(logDir, rl.baseName+"-*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackupFiles {
		return nil
	}

	filesToRemove := matches[:len(matches)-maxBackupFiles]
	for _, file := range filesToRemove {
		compressedFile := file + ".gz"
		if err := compressFile(file, compressedFile); err != nil {
			log.Printf("Compression failed for %s: %v", file, err)
		}
		if err := os.Remove(file); err != nil {
			log.Printf("Failed to remove %s: %v", file, err)
		}
	}

	return nil
}

func compressFile(src, dst string) error {
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

	gzWriter := io.MultiWriter(destination)
	_, err = io.Copy(gzWriter, source)
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "ROTATING: ", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d with some sample data", i)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}
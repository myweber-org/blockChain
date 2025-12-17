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
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:  basePath,
		maxSize:   int64(maxSizeMB) * 1024 * 1024,
		fileIndex: 0,
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

	return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentFile != nil {
		info, err := rl.currentFile.Stat()
		if err != nil {
			return err
		}
		if info.Size() < rl.maxSize {
			return nil
		}
		rl.currentFile.Close()
		if err := rl.compressOldLog(); err != nil {
			fmt.Printf("Failed to compress log: %v\n", err)
		}
	}

	filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileIndex)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.currentFile = f
	rl.fileIndex++
	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	if rl.fileIndex == 0 {
		return nil
	}
	oldIndex := rl.fileIndex - 1
	srcName := fmt.Sprintf("%s.%d.log", rl.basePath, oldIndex)
	dstName := fmt.Sprintf("%s.%d.log.gz", rl.basePath, oldIndex)

	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	if err := os.Remove(srcName); err != nil {
		return err
	}
	return nil
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
	logger, err := NewRotatingLogger("app", 10)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Processing request from client\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	files, _ := filepath.Glob("app.*")
	fmt.Printf("Generated files: %v\n", files)
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	baseName    string
}

func NewLogRotator(name string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	lr := &LogRotator{baseName: basePath}

	if err := lr.openCurrent(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrent() error {
	path := lr.baseName + ".log"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102-150405")
	oldPath := fmt.Sprintf("%s-%s.log", lr.baseName, timestamp)
	if err := os.Rename(lr.baseName+".log", oldPath); err != nil {
		return err
	}

	if err := lr.openCurrent(); err != nil {
		return err
	}

	lr.cleanupOld()
	return nil
}

func (lr *LogRotator) cleanupOld() {
	pattern := lr.baseName + "-*.log"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	toDelete := matches[:len(matches)-maxBackups]
	for _, path := range toDelete {
		os.Remove(path)
	}
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	writer := io.MultiWriter(os.Stdout, rotator)
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		writer.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}
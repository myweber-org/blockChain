
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
	maxFileSize = 10 * 1024 * 1024
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	baseName  string
	fileSize  int64
	fileIndex int
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
	rl.fileIndex = rl.findLatestIndex()
	filename := rl.generateFilename(rl.fileIndex)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) findLatestIndex() int {
	pattern := filepath.Join(logDir, rl.baseName+"*.log")
	matches, _ := filepath.Glob(pattern)
	maxIndex := 0

	for _, match := range matches {
		var index int
		fmt.Sscanf(filepath.Base(match), rl.baseName+".%d.log", &index)
		if index > maxIndex {
			maxIndex = index
		}
	}

	return maxIndex
}

func (rl *RotatingLogger) generateFilename(index int) string {
	return filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, index))
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.fileSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
		rl.compressCurrentFile()
	}

	rl.fileIndex++
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressCurrentFile() {
	oldFile := rl.generateFilename(rl.fileIndex)
	compressedFile := oldFile + ".gz"

	go func() {
		src, err := os.Open(oldFile)
		if err != nil {
			return
		}
		defer src.Close()

		dst, err := os.Create(compressedFile)
		if err != nil {
			return
		}
		defer dst.Close()

		gz := gzip.NewWriter(dst)
		defer gz.Close()

		if _, err := io.Copy(gz, src); err != nil {
			return
		}

		os.Remove(oldFile)
	}()
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
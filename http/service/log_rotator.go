
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
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
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
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.file != nil && rl.size < maxFileSize {
		return nil
	}

	if rl.file != nil {
		rl.file.Close()
		if err := rl.compressCurrent(); err != nil {
			return err
		}
		rl.cleanOldBackups()
	}

	newFile, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := newFile.Stat()
	if err != nil {
		newFile.Close()
		return err
	}

	rl.file = newFile
	rl.size = info.Size()
	rl.currentNum = 0
	return nil
}

func (rl *RotatingLogger) compressCurrent() error {
	src, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer src.Close()

	backupName := fmt.Sprintf("%s.%s.gz", rl.basePath, time.Now().Format("20060102_150405"))
	dst, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	_, err = io.Copy(gz, src)
	return err
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
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
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))))
		time.Sleep(10 * time.Millisecond)
	}
}
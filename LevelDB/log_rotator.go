
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	basePath      string
	maxSize       int64
	currentSize   int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		currentFile: file,
		basePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	rl.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(rl.basePath)
	base := strings.TrimSuffix(rl.basePath, ext)
	archivePath := fmt.Sprintf("%s_%s_%d%s", base, timestamp, rl.rotationCount, ext)

	if err := os.Rename(rl.basePath, archivePath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(io.MultiWriter(os.Stdout, logger), "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}
}
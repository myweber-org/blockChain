
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		backupCount: backupCount,
		compressOld: compressOld,
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

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
			log.Printf("Failed to rotate log: %v", err)
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

	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		
		if i == 0 {
			oldPath = rl.basePath
		} else {
			oldPath = filepath.Join(dir, fmt.Sprintf("%s.%d", baseName, i))
			if rl.compressOld {
				oldPath += ".gz"
			}
		}

		newPath = filepath.Join(dir, fmt.Sprintf("%s.%d", baseName, i+1))
		if rl.compressOld && i+1 < rl.backupCount {
			newPath += ".gz"
		}

		if _, err := os.Stat(oldPath); err == nil {
			if i+1 >= rl.backupCount {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if rl.compressOld {
		go rl.compressFile(rl.basePath, rl.basePath+".1.gz")
	} else {
		os.Rename(rl.basePath, rl.basePath+".1")
	}

	return rl.openOrCreate()
}

func (rl *RotatingLogger) compressFile(src, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	compressor := NewGzipCompressor(dstFile)
	_, err = io.Copy(compressor, srcFile)
	compressor.Close()

	if err == nil {
		os.Remove(src)
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

type GzipCompressor struct {
	w io.WriteCloser
}

func NewGzipCompressor(w io.Writer) *GzipCompressor {
	return &GzipCompressor{}
}

func (g *GzipCompressor) Write(p []byte) (int, error) {
	return len(p), nil
}

func (g *GzipCompressor) Close() error {
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: Application event at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}
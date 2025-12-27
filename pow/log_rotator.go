
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingWriter struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingWriter(filePath string, maxSize int64) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := writer.openCurrentFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (rw *RotatingWriter) openCurrentFile() error {
	dir := filepath.Dir(rw.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rw.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rw.currentFile = file
	rw.currentSize = info.Size()
	return nil
}

func (rw *RotatingWriter) rotate() error {
	if rw.currentFile != nil {
		rw.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", rw.filePath, timestamp, rw.rotationCount)
	rw.rotationCount++

	if err := os.Rename(rw.filePath, backupPath); err != nil {
		return err
	}

	return rw.openCurrentFile()
}

func (rw *RotatingWriter) Write(p []byte) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.currentSize+int64(len(p)) > rw.maxSize {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rw.currentFile.Write(p)
	if err == nil {
		rw.currentSize += int64(n)
	}
	return n, err
}

func (rw *RotatingWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.currentFile != nil {
		return rw.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create rotating writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		_, err := writer.Write([]byte(message))
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
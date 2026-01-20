
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
	current      *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
	w := &RotatingWriter{
		basePath: basePath,
		maxSize:  maxSize,
	}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%d_%s.log", w.basePath, w.rotationCount, timestamp)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.current = file
	info, err := file.Stat()
	if err == nil {
		w.currentSize = info.Size()
	}
	return nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		w.current.Close()
		w.rotationCount++
		if err := w.openFile(); err != nil {
			return 0, err
		}
	}

	n, err := w.current.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.current.Close()
}

func main() {
	writer, err := NewRotatingWriter("app_log", 1024*1024)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		writer.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}
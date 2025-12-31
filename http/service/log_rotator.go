package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type RotatingFile struct {
	mu          sync.Mutex
	file        *os.File
	basePath    string
	maxSize     int64
	maxFiles    int
	currentSize int64
}

func NewRotatingFile(basePath string, maxSize int64, maxFiles int) (*RotatingFile, error) {
	rf := &RotatingFile{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rf.openCurrentFile(); err != nil {
		return nil, err
	}

	return rf, nil
}

func (rf *RotatingFile) openCurrentFile() error {
	file, err := os.OpenFile(rf.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rf.file = file
	rf.currentSize = info.Size()
	return nil
}

func (rf *RotatingFile) rotate() error {
	rf.file.Close()

	for i := rf.maxFiles - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", rf.basePath, i)
		newPath := fmt.Sprintf("%s.%d", rf.basePath, i+1)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	backupPath := fmt.Sprintf("%s.1", rf.basePath)
	os.Rename(rf.basePath, backupPath)

	return rf.openCurrentFile()
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.currentSize+int64(len(p)) > rf.maxSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.file.Write(p)
	if err == nil {
		rf.currentSize += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) Close() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.file.Close()
}

func main() {
	logFile, err := NewRotatingFile("app.log", 1024*1024, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a test log message.\n", i)
		if _, err := logFile.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
	}

	fmt.Println("Log rotation test completed. Check app.log and rotated files.")
}
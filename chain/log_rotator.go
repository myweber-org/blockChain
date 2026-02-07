
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
	filename    string
	maxSize     int64
	currentSize int64
	file        *os.File
	rotation    int
}

func NewRotatingFile(filename string, maxSize int64) (*RotatingFile, error) {
	rf := &RotatingFile{
		filename: filename,
		maxSize:  maxSize,
	}

	if err := rf.openFile(); err != nil {
		return nil, err
	}

	return rf, nil
}

func (rf *RotatingFile) openFile() error {
	info, err := os.Stat(rf.filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if info != nil {
		rf.currentSize = info.Size()
	}

	file, err := os.OpenFile(rf.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rf.file = file
	return nil
}

func (rf *RotatingFile) rotate() error {
	if rf.file != nil {
		rf.file.Close()
	}

	backupName := fmt.Sprintf("%s.%d", rf.filename, rf.rotation)
	rf.rotation++

	if err := os.Rename(rf.filename, backupName); err != nil {
		return err
	}

	rf.currentSize = 0
	return rf.openFile()
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

	if rf.file != nil {
		return rf.file.Close()
	}
	return nil
}

func main() {
	logFile, err := NewRotatingFile("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}
	defer logFile.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
		if _, err := logFile.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
	}

	fmt.Println("Log rotation test completed")
}
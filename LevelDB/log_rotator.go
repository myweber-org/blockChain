
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
	backupCount = 5
)

type RotatingFile struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingFile(path string) (*RotatingFile, error) {
	rf := &RotatingFile{
		basePath: path,
	}
	if err := rf.openCurrent(); err != nil {
		return nil, err
	}
	return rf, nil
}

func (rf *RotatingFile) openCurrent() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.file != nil {
		rf.file.Close()
	}

	file, err := os.OpenFile(rf.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rf.file = file
	rf.size = stat.Size()
	return nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.size+int64(len(p)) > maxFileSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.file.Write(p)
	if err == nil {
		rf.size += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotate() error {
	if rf.file != nil {
		rf.file.Close()
		rf.file = nil
	}

	for i := backupCount - 1; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d.gz", rf.basePath, i)
		newName := fmt.Sprintf("%s.%d.gz", rf.basePath, i+1)
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	if err := rf.compressCurrent(); err != nil {
		return err
	}

	return rf.openCurrent()
}

func (rf *RotatingFile) compressCurrent() error {
	src, err := os.Open(rf.basePath)
	if err != nil {
		return err
	}
	defer src.Close()

	destPath := fmt.Sprintf("%s.1.gz", rf.basePath)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(rf.basePath)
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
	logFile, err := NewRotatingFile("application.log")
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}
	defer logFile.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Processing request from client\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logFile.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
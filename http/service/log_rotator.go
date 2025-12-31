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
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
	maxRotations int
}

func NewRotatingLogger(filePath string, maxSizeMB int, maxRotations int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		currentFile:  file,
		filePath:     filePath,
		maxSize:      maxSize,
		currentSize:  info.Size(),
		rotationCount: 0,
		maxRotations: maxRotations,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

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

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, rotatedPath); err != nil {
		return err
	}

	if err := rl.compressFile(rotatedPath); err != nil {
		fmt.Printf("Failed to compress %s: %v\n", rotatedPath, err)
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	rl.rotationCount++

	if rl.rotationCount > rl.maxRotations {
		rl.cleanupOldLogs()
	}

	return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(source + ".gz")
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		return err
	}

	if err := os.Remove(source); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
	pattern := rl.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxRotations {
		filesToDelete := matches[:len(matches)-rl.maxRotations]
		for _, file := range filesToDelete {
			os.Remove(file)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}

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
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
	currentSize int64
}

func NewRotatingLogger(filePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	logger := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	l.currentFile = file
	l.currentSize = info.Size()
	return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentSize+int64(len(p)) > l.maxSize {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := l.currentFile.Write(p)
	if err == nil {
		l.currentSize += int64(n)
	}
	return n, err
}

func (l *RotatingLogger) rotate() error {
	if err := l.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", l.filePath, timestamp)

	if err := os.Rename(l.filePath, backupPath); err != nil {
		return err
	}

	if err := l.compressFile(backupPath); err != nil {
		return err
	}

	if err := l.cleanupOldBackups(); err != nil {
		return err
	}

	return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	compressedPath := sourcePath + ".gz"
	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	gzipWriter := gzip.NewWriter(compressedFile)
	defer gzipWriter.Close()

	if _, err := io.Copy(gzipWriter, sourceFile); err != nil {
		return err
	}

	if err := os.Remove(sourcePath); err != nil {
		return err
	}

	return nil
}

func (l *RotatingLogger) cleanupOldBackups() error {
	pattern := l.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= l.backupCount {
		return nil
	}

	oldestFiles := matches[:len(matches)-l.backupCount]
	for _, file := range oldestFiles {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (l *RotatingLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentFile != nil {
		return l.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
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
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentSize < rl.maxSize {
		return nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.currentFile.Close()
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationCount)
	rl.rotationCount++

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) WriteString(s string) (int, error) {
	return rl.Write([]byte(s))
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.WriteString(message); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
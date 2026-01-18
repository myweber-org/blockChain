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
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		basePath: basePath,
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

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.currentFile = nil
	}

	rl.rotationCount++
	archiveName := fmt.Sprintf("%s.%d-%s.gz", 
		rl.basePath, 
		rl.rotationCount,
		time.Now().Format("20060102-150405"))

	source, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer dest.Close()

	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, source); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	return rl.openCurrentFile()
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
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
	fileIndex   int
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rotator := &LogRotator{
		filePath:  filePath,
		maxSize:   maxSize,
		fileIndex: 0,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	if lr.file != nil {
		lr.file.Close()
	}

	var err error
	lr.file, err = os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := lr.file.Stat()
	if err != nil {
		return err
	}

	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) rotateIfNeeded() error {
	if lr.currentSize >= lr.maxSize {
		lr.file.Close()

		backupPath := lr.filePath + "." + strconv.Itoa(lr.fileIndex) + "." + time.Now().Format("20060102")
		err := os.Rename(lr.filePath, backupPath)
		if err != nil {
			return err
		}

		lr.fileIndex++
		return lr.openCurrentFile()
	}
	return nil
}

func (lr *LogRotator) WriteLog(level, message string) error {
	if err := lr.rotateIfNeeded(); err != nil {
		return err
	}

	logEntry := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, message)
	n, err := lr.file.WriteString(logEntry)
	if err != nil {
		return err
	}

	lr.currentSize += int64(n)
	return nil
}

func (lr *LogRotator) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("application.log", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		err := rotator.WriteLog("INFO", fmt.Sprintf("Log entry number %d", i))
		if err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
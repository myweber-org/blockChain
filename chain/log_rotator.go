
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
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	mu          sync.Mutex
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
	}

	timestamp := time.Now().Format("20060102-150405")
	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressFile(newPath); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := filepath.Join(logDir, rl.baseName+"-*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		filesToRemove := matches[:len(matches)-maxBackups]
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}

	return nil
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const maxLogSize = 10 * 1024 * 1024 // 10MB

type LogRotator struct {
	filePath string
	currentSize int64
}

func NewLogRotator(path string) *LogRotator {
	return &LogRotator{
		filePath: path,
	}
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.shouldRotate(int64(len(p))) {
		err := lr.rotate()
		if err != nil {
			return 0, err
		}
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err := file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) shouldRotate(additionalBytes int64) bool {
	info, err := os.Stat(lr.filePath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.Size()+additionalBytes > maxLogSize
}

func (lr *LogRotator) rotate() error {
	timestamp := time.Now().Unix()
	backupPath := lr.filePath + "." + strconv.FormatInt(timestamp, 10)
	
	err := os.Rename(lr.filePath, backupPath)
	if err != nil {
		return err
	}
	
	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) CleanOldLogs(maxBackups int) error {
	pattern := lr.filePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	
	if len(matches) <= maxBackups {
		return nil
	}
	
	for i := 0; i < len(matches)-maxBackups; i++ {
		err := os.Remove(matches[i])
		if err != nil {
			return err
		}
	}
	
	return nil
}
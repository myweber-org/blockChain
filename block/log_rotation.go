package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	maxFileSize = 1024 * 1024
	maxBackups  = 5
)

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	filePath string
	size     int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		file:     file,
		filePath: path,
		size:     info.Size(),
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	for i := maxBackups - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(rl.filePath, rl.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.size = 0
	return nil
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d", rl.filePath, index+1)
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: This is a test log message\n", i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"log"
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
	mu        sync.Mutex
	file      *os.File
	size      int64
	baseName  string
	fileIndex int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	filename := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = info.Size()
	rl.fileIndex = 0
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	rl.fileIndex++
	if rl.fileIndex > maxBackups {
		rl.fileIndex = maxBackups
	}

	for i := rl.fileIndex; i > 0; i-- {
		oldName := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, i-1))
		newName := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, i))
		if i == maxBackups {
			os.Remove(newName)
		}
		os.Rename(oldName, newName)
	}

	os.Rename(
		filepath.Join(logDir, rl.baseName+".log"),
		filepath.Join(logDir, rl.baseName+".0.log"),
	)

	return rl.openCurrent()
}

func (rl *RotatingLogger) CleanupOld() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for i := maxBackups; i < 20; i++ {
		filename := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, i))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		os.Remove(filename)
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			logger.CleanupOld()
		}
	}()

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application event occurred at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}
}
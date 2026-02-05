
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        sequence: 0,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        lr.file.Close()
    }

    filename := lr.basePath
    if lr.sequence > 0 {
        filename = fmt.Sprintf("%s.%d", lr.basePath, lr.sequence)
    }

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.file = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) rotate() error {
    lr.sequence++
    return lr.openCurrentFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func (lr *LogRotator) CleanOldLogs(maxFiles int) error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    for i := lr.sequence - maxFiles; i >= 0; i-- {
        filename := lr.basePath
        if i > 0 {
            filename = fmt.Sprintf("%s.%d", lr.basePath, i)
        }
        if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
            return err
        }
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message.\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    if err := rotator.CleanOldLogs(5); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }
}package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	basePath   string
	maxSize    int64
	currentSize int64
	backupCount int
}

func NewRotatingLogger(basePath string, maxSize int64, backupCount int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:   basePath,
		maxSize:    maxSize,
		backupCount: backupCount,
	}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}
	for i := rl.backupCount - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}
	if err := os.Rename(rl.basePath, rl.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, index)
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10*1024*1024, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			log.Printf("Write error: %v", err)
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
)

type LogRotator struct {
	mu          sync.Mutex
	file        *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	backupCount int
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) (*LogRotator, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &LogRotator{
		file:        file,
		filePath:    filePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		backupCount: backupCount,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	n, err := lr.file.Write(p)
	if err != nil {
		return n, fmt.Errorf("failed to write to log: %w", err)
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	for i := lr.backupCount - 1; i >= 0; i-- {
		oldPath := lr.backupPath(i)
		newPath := lr.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rename backup file: %w", err)
			}
		}
	}

	if err := os.Rename(lr.filePath, lr.backupPath(0)); err != nil {
		return fmt.Errorf("failed to rename current log file: %w", err)
	}

	file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	lr.file = file
	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) backupPath(index int) string {
	if index == 0 {
		return lr.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d", lr.filePath, index+1)
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i+1)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize   = 1024 * 1024 // 1MB
	logDirectory = "./logs"
	baseLogName  = "app.log"
)

func rotateLogIfNeeded() error {
	logPath := filepath.Join(logDirectory, baseLogName)
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if info.Size() < maxLogSize {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := fmt.Sprintf("app_%s.log", timestamp)
	archivePath := filepath.Join(logDirectory, archiveName)

	err = os.Rename(logPath, archivePath)
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	fmt.Printf("Log rotated: %s -> %s\n", baseLogName, archiveName)
	return nil
}

func ensureLogDirectory() error {
	return os.MkdirAll(logDirectory, 0755)
}

func writeLogEntry(message string) error {
	err := ensureLogDirectory()
	if err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	err = rotateLogIfNeeded()
	if err != nil {
		return fmt.Errorf("failed to rotate log: %w", err)
	}

	logPath := filepath.Join(logDirectory, baseLogName)
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err = file.WriteString(logLine)
	if err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

func main() {
	for i := 1; i <= 100; i++ {
		message := fmt.Sprintf("Test log entry number %d", i)
		err := writeLogEntry(message)
		if err != nil {
			fmt.Printf("Error writing log: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("Log rotation test completed")
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLog struct {
    mu          sync.Mutex
    file        *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    maxFiles    int
}

func NewRotatingLog(basePath string, maxSizeMB int, maxFiles int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        file:        file,
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        maxFiles:    maxFiles,
    }, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    compressedPath := rotatedPath + ".gz"
    if err := compressFile(rotatedPath, compressedPath); err != nil {
        return err
    }
    os.Remove(rotatedPath)

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0

    go rl.cleanupOldFiles()

    return nil
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (rl *RotatingLog) cleanupOldFiles() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > rl.maxFiles {
        filesToDelete := matches[:len(matches)-rl.maxFiles]
        for _, f := range filesToDelete {
            os.Remove(f)
        }
    }
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    log, err := NewRotatingLog("app.log", 10, 5)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        log.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().String())))
        time.Sleep(10 * time.Millisecond)
    }
}
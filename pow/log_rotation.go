
package main

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	file        *os.File
	basePath    string
	maxSize     int64
	currentSize int64
	fileCount   int
	maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
		maxFiles: maxFiles,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
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

	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentSize < rl.maxSize {
		return nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	for i := rl.maxFiles - 1; i > 0; i-- {
		oldPath := rl.basePath + "." + strconv.Itoa(i)
		newPath := rl.basePath + "." + strconv.Itoa(i+1)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	backupPath := rl.basePath + ".1"
	os.Rename(rl.basePath, backupPath)

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.file.Write(p)
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

	if rl.file != nil {
		return rl.file.Close()
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
		logger.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - Log entry " + strconv.Itoa(i) + "\n")
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxLogSize   = 10 * 1024 * 1024 // 10MB
	backupCount  = 5
	logDir       = "./logs"
	currentLog   = "app.log"
)

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	size     int64
	basePath string
}

func NewRotatingLogger() (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	path := filepath.Join(logDir, currentLog)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	info, _ := file.Stat()
	return &RotatingLogger{
		file:     file,
		size:     info.Size(),
		basePath: path,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.size += int64(n)
	if rl.size >= maxLogSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()

	backupPath := fmt.Sprintf("%s.%s.gz",
		rl.basePath,
		time.Now().Format("20060102-150405"),
	)

	src, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.size = 0
	rl.cleanOldBackups()
	return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := filepath.Join(logDir, currentLog+".*.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	toDelete := matches[:len(matches)-backupCount]
	for _, path := range toDelete {
		os.Remove(path)
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 10000; i++ {
		log.Printf("Log entry %d: Application event processed", i)
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 1024 * 1024 // 1MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    fileSize    int64
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    fullPath := filepath.Join(path, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile: file,
        fileSize:    info.Size(),
        basePath:    path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.fileSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.fileSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(rl.basePath, logFileName)
    newPath := filepath.Join(rl.basePath, backupName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.fileSize = 0

    go rl.cleanupOldFiles()
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := filepath.Join(rl.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackupFiles {
        filesToDelete := matches[:len(matches)-maxBackupFiles]
        for _, f := range filesToDelete {
            os.Remove(f)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }
}
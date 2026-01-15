package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentSize int64
	maxSize     int64
	basePath    string
	file        *os.File
	fileIndex   int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:  basePath,
		maxSize:   maxSize,
		fileIndex: 0,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := rl.generateFilename()
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) generateFilename() string {
	if rl.fileIndex == 0 {
		return rl.basePath
	}
	ext := filepath.Ext(rl.basePath)
	base := rl.basePath[:len(rl.basePath)-len(ext)]
	return fmt.Sprintf("%s.%d%s", base, rl.fileIndex, ext)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	rl.fileIndex++
	return rl.openCurrentFile()
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
	logger, err := NewRotatingLogger("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message for testing rotation.\n", i)
		logger.Write([]byte(message))
	}

	fmt.Println("Log rotation test completed")
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu           sync.Mutex
    basePath     string
    maxSize      int64
    currentSize  int64
    currentFile  *os.File
    sequence     int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        sequence: 0,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
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
    rl.sequence = rl.findLatestSequence()
    return nil
}

func (rl *RotatingLogger) findLatestSequence() int {
    pattern := rl.basePath + ".*.gz"
    matches, _ := filepath.Glob(pattern)
    maxSeq := 0
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        seq, err := strconv.Atoi(parts[len(parts)-2])
        if err == nil && seq > maxSeq {
            maxSeq = seq
        }
    }
    return maxSeq
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
    }

    timestamp := time.Now().Format("20060102_150405")
    rl.sequence++
    archiveName := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.sequence, timestamp)

    if err := rl.compressFile(rl.basePath, archiveName); err != nil {
        return err
    }

    if err := os.Truncate(rl.basePath, 0); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
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
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application event processed\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)

type LogRotator struct {
    filePath    string
    maxSize     int64
    backupCount int
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) *LogRotator {
    return &LogRotator{
        filePath:    filePath,
        maxSize:     maxSize,
        backupCount: backupCount,
    }
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if err := lr.rotateIfNeeded(); err != nil {
        return 0, err
    }

    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    return file.Write(p)
}

func (lr *LogRotator) rotateIfNeeded() error {
    info, err := os.Stat(lr.filePath)
    if os.IsNotExist(err) {
        return nil
    }
    if err != nil {
        return err
    }

    if info.Size() < lr.maxSize {
        return nil
    }

    for i := lr.backupCount - 1; i >= 0; i-- {
        oldPath := lr.backupPath(i)
        newPath := lr.backupPath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            if err := os.Rename(oldPath, newPath); err != nil {
                return err
            }
        }
    }

    if err := os.Rename(lr.filePath, lr.backupPath(0)); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) backupPath(index int) string {
    if index == 0 {
        return lr.filePath + ".1"
    }
    return fmt.Sprintf("%s.%d", lr.filePath, index+1)
}

func main() {
    rotator := NewLogRotator("app.log", 1024*1024, 5)

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Error writing log: %v\n", err)
        }
    }

    fmt.Println("Log rotation test completed")
}
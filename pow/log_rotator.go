
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile   *os.File
    currentSize   int64
    basePath      string
    currentIndex  int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        err := lr.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err != nil {
        return n, err
    }

    lr.currentSize += int64(n)
    return n, nil
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.gz", lr.basePath, timestamp)

    err := compressFile(lr.basePath, rotatedPath)
    if err != nil {
        return err
    }

    os.Remove(lr.basePath)

    err = lr.cleanupOldBackups()
    if err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    backupTimes := make([]time.Time, len(matches))
    for i, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", timestamp)
        if err != nil {
            continue
        }
        backupTimes[i] = t
    }

    oldestIndex := 0
    for i := 1; i < len(backupTimes); i++ {
        if backupTimes[i].Before(backupTimes[oldestIndex]) {
            oldestIndex = i
        }
    }

    if oldestIndex < len(matches) {
        return os.Remove(matches[oldestIndex])
    }

    return nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
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
	maxFileSize = 10 * 1024 * 1024
	maxFiles    = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	fileSize    int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
		sequence: 0,
	}

	if err := rl.openNextFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	rl.sequence++
	if rl.sequence > maxFiles {
		rl.sequence = 1
	}

	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.fileSize = 0

	rl.cleanOldFiles()
	return nil
}

func (rl *RotatingLogger) cleanOldFiles() {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if time.Since(info.ModTime()) > 7*24*time.Hour {
			os.Remove(file)
		}
	}
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.fileSize+int64(len(p)) > maxFileSize {
		if err := rl.openNextFile(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	rl.fileSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
		time.Sleep(100 * time.Millisecond)
	}
}

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

type RotatingLog struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
}

func NewRotatingLog(path string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        filePath:    path,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
        sequence:    0,
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
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archiveName := fmt.Sprintf("%s.%s.%d.gz", rl.filePath, timestamp, rl.sequence)
    rl.sequence++

    if err := rl.compressFile(rl.filePath, archiveName); err != nil {
        return err
    }

    if err := os.Truncate(rl.filePath, 0); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func (rl *RotatingLog) compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logFile := "application.log"
    rotator, err := NewRotatingLog(logFile, 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
    files, _ := filepath.Glob("*.gz")
    fmt.Printf("Created %d compressed archives\n", len(files))
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

type LogRotator struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
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

	return &LogRotator{
		currentFile:  file,
		filePath:     filePath,
		maxSize:      maxSize,
		currentSize:  info.Size(),
		rotationCount: 0,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%d.%s.gz", lr.filePath, lr.rotationCount, timestamp)
	lr.rotationCount++

	if err := compressFile(lr.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Truncate(lr.filePath, 0); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0
	return nil
}

func compressFile(source, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	gzWriter := gzip.NewWriter(targetFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, sourceFile)
	return err
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
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
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
package main

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

type RotatingLogger struct {
    mu          sync.Mutex
    file        *os.File
    currentSize int64
    maxSize     int64
    basePath    string
    sequence    int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        maxSize:  maxSize,
        basePath: basePath,
        sequence: 0,
    }
    if err := rl.openLogFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openLogFile() error {
    filename := fmt.Sprintf("%s.log", rl.basePath)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            log.Printf("Failed to rotate log: %v", err)
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

    timestamp := time.Now().Format("20060102_150405")
    oldPath := fmt.Sprintf("%s_%s_%03d.log", rl.basePath, timestamp, rl.sequence)
    rl.sequence++

    if err := os.Rename(fmt.Sprintf("%s.log", rl.basePath), oldPath); err != nil {
        return err
    }

    if err := rl.compressFile(oldPath); err != nil {
        log.Printf("Failed to compress %s: %v", oldPath, err)
    }

    return rl.openLogFile()
}

func (rl *RotatingLogger) compressFile(source string) error {
    reader, err := os.Open(source)
    if err != nil {
        return err
    }
    defer reader.Close()

    compressed := source + ".gz"
    writer, err := os.Create(compressed)
    if err != nil {
        return err
    }
    defer writer.Close()

    gz := gzip.NewWriter(writer)
    defer gz.Close()

    if _, err := io.Copy(gz, reader); err != nil {
        return err
    }

    if err := os.Remove(source); err != nil {
        return err
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
    logger, err := NewRotatingLogger("app", 1024*1024)
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }

    files, _ := filepath.Glob("app_*.log.gz")
    fmt.Printf("Created %d compressed log files\n", len(files))
}
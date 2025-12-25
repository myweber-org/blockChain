
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
    mu            sync.Mutex
    currentFile   *os.File
    basePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
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

    return &LogRotator{
        currentFile: file,
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
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
    archivedPath := fmt.Sprintf("%s.%s.gz", lr.basePath, timestamp)

    sourceFile, err := os.Open(lr.basePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(archivedPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(lr.basePath)

    newFile, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = newFile
    lr.currentSize = 0
    lr.rotationCount++

    return nil
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
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
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
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu            sync.Mutex
    basePath      string
    currentFile   *os.File
    maxSize       int64
    currentSize   int64
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

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)

    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }

    l.rotationCount++
    go l.compressOldLog(rotatedPath)

    return l.openCurrentFile()
}

func (l *RotatingLogger) compressOldLog(path string) {
    compressedPath := path + ".gz"
    src, err := os.Open(path)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return
    }

    os.Remove(path)
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func findLatestRotation(basePath string) int {
    pattern := basePath + ".*"
    matches, _ := filepath.Glob(pattern)
    maxNum := 0

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }
        if num, err := strconv.Atoi(parts[len(parts)-1]); err == nil && num > maxNum {
            maxNum = num
        }
    }
    return maxNum
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
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println("Log rotation example completed")
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
    rotator := &LogRotator{
        filePath: filePath,
        maxSize:  maxSize,
    }

    if err := rotator.openFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openFile() error {
    if lr.file != nil {
        lr.file.Close()
    }

    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    return lr.openFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
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
    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func main() {
    logDir := "./logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        fmt.Printf("Failed to create log directory: %v\n", err)
        return
    }

    logPath := filepath.Join(logDir, "application.log")
    rotator, err := NewLogRotator(logPath, 1024*1024) // 1MB max size
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
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
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    file        *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
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

    return &RotatingLogger{
        file:        file,
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        maxFiles:    maxFiles,
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

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.cleanupOldFiles()

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

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    os.Remove(source)
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= rl.maxFiles {
        return
    }

    filesToRemove := matches[:len(matches)-rl.maxFiles]
    for _, file := range filesToRemove {
        os.Remove(file)
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	filePath     string
	maxSize      int64
	currentSize  int64
	file         *os.File
	rotationCount int
}

func NewRotatingLogger(path string, maxSizeMB int) (*RotatingLogger, error) {
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
	
	return &RotatingLogger{
		filePath:    path,
		maxSize:     maxSize,
		currentSize: info.Size(),
		file:        file,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
	rl.file.Close()
	
	rl.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(rl.filePath)
	base := strings.TrimSuffix(rl.filePath, ext)
	archivePath := fmt.Sprintf("%s_%s_%d%s", base, timestamp, rl.rotationCount, ext)
	
	if err := os.Rename(rl.filePath, archivePath); err != nil {
		return err
	}
	
	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	rl.file = file
	rl.currentSize = 0
	
	go rl.compressArchive(archivePath)
	
	return nil
}

func (rl *RotatingLogger) compressArchive(path string) {
	src, err := os.Open(path)
	if err != nil {
		return
	}
	defer src.Close()
	
	destPath := path + ".gz"
	dest, err := os.Create(destPath)
	if err != nil {
		return
	}
	defer dest.Close()
	
	gzWriter := NewGzipWriter(dest)
	defer gzWriter.Close()
	
	io.Copy(gzWriter, src)
	
	os.Remove(path)
}

func (rl *RotatingLogger) Close() error {
	return rl.file.Close()
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
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    if err := os.MkdirAll(basePath, 0755); err != nil {
        return nil, err
    }

    fullPath := filepath.Join(basePath, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        currentFile: file,
        currentSize: info.Size(),
        basePath:    basePath,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = lr.currentFile.Write(p)
    lr.currentSize += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(lr.basePath, logFileName)
    newPath := filepath.Join(lr.basePath, backupName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    lr.currentSize = 0

    go lr.cleanupOldLogs()
    return nil
}

func (lr *LogRotator) cleanupOldLogs() {
    files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i, file := range files {
        if i >= maxBackupCount {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) Close() error {
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Test data for rotation check\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
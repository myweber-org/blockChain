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
	basePath    string
	maxSize     int64
	currentSize int64
	fileCount   int
	maxFiles    int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
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

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.currentFile = f

	info, err := f.Stat()
	if err != nil {
		return err
	}
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
	rl.currentFile.Close()

	oldPath := fmt.Sprintf("%s.%d", rl.basePath, rl.fileCount)
	if err := os.Rename(rl.basePath, oldPath); err != nil {
		return err
	}

	if err := rl.compressFile(oldPath); err != nil {
		return err
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		oldest := rl.fileCount - rl.maxFiles - 1
		oldFile := fmt.Sprintf("%s.%d.gz", rl.basePath, oldest)
		os.Remove(oldFile)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	defer os.Remove(src)

	dst := src + ".gz"
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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
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
    "strconv"
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    currentSize   int64
    currentFile   *os.File
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
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        currentFile: file,
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
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    lr.currentSize = 0
    lr.rotationCount++

    lr.cleanupOldLogs()
    return nil
}

func (lr *LogRotator) compressFile(sourcePath string) error {
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

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldLogs() {
    dir := filepath.Dir(lr.basePath)
    baseName := filepath.Base(lr.basePath)

    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var compressedLogs []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedLogs = append(compressedLogs, filepath.Join(dir, name))
        }
    }

    if len(compressedLogs) > 10 {
        for i := 0; i < len(compressedLogs)-10; i++ {
            os.Remove(compressedLogs[i])
        }
    }
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
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))

        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	mu          sync.Mutex
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
	fileCount   int
	maxFiles    int
}

func NewLogRotator(filePath string, maxSize int64, maxFiles int) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
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
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.file != nil {
		lr.file.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, rotatedPath); err != nil {
		return err
	}

	lr.fileCount++

	if lr.fileCount > lr.maxFiles {
		lr.cleanOldFiles()
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) cleanOldFiles() {
	pattern := lr.filePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > lr.maxFiles {
		filesToRemove := len(matches) - lr.maxFiles
		for i := 0; i < filesToRemove; i++ {
			os.Remove(matches[i])
		}
	}
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		lr.mu.Unlock()
		if err := lr.rotate(); err != nil {
			lr.mu.Lock()
			return 0, err
		}
		lr.mu.Lock()
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

func main() {
	rotator, err := NewLogRotator("app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		rotator.Write([]byte(message))
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

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
			log.Printf("Failed to rotate log: %v", err)
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
		log.Printf("Failed to compress %s: %v", newPath, err)
	}

	if err := rl.cleanupOldFiles(); err != nil {
		log.Printf("Failed to cleanup old files: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	dstPath := path + ".gz"
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
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
	logger, err := NewRotatingLogger("application")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, "Sample log message for testing rotation")
		time.Sleep(10 * time.Millisecond)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	basePath      string
	maxSize       int64
	currentSize   int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
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
		currentFile: file,
		basePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	rl.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(rl.basePath)
	base := strings.TrimSuffix(rl.basePath, ext)
	archivePath := fmt.Sprintf("%s_%s_%d%s", base, timestamp, rl.rotationCount, ext)

	if err := os.Rename(rl.basePath, archivePath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(io.MultiWriter(os.Stdout, logger), "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: path}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	rl.currentSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := rl.basePath
	newPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	if err := compressFile(oldPath, newPath); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	cleanupOldBackups(rl.basePath)
	return rl.openCurrent()
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

func cleanupOldBackups(basePath string) {
	pattern := basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, f := range toDelete {
			os.Remove(f)
		}
	}
}

func (rl *RotatingLogger) openCurrent() error {
	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().String())))
		time.Sleep(10 * time.Millisecond)
	}
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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
	currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, backups int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	logger := &RotatingLogger{
		filePath:    absPath,
		maxSize:     maxSize,
		backupCount: backups,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	if err := rl.compressBackup(backupPath); err != nil {
		return err
	}

	if err := rl.cleanOldBackups(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressBackup(srcPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := srcPath + ".gz"
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(srcPath); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := rl.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= rl.backupCount {
		return nil
	}

	toDelete := len(matches) - rl.backupCount
	for i := 0; i < toDelete; i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
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
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	compressOld = true
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDate {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLogs(); err != nil {
				log.Printf("Failed to compress old logs: %v", err)
			}
		}

		newPath := rl.generateLogPath(now)
		if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
		rl.currentDay = currentDate
	}
	return nil
}

func (rl *RotatingLogger) generateLogPath(t time.Time) string {
	base := strings.TrimSuffix(rl.basePath, filepath.Ext(rl.basePath))
	ext := filepath.Ext(rl.basePath)
	if ext == "" {
		ext = ".log"
	}

	for i := 0; i < 100; i++ {
		suffix := ""
		if i > 0 {
			suffix = fmt.Sprintf("-%d", i)
		}
		path := fmt.Sprintf("%s-%s%s%s", base, t.Format("2006-01-02"), suffix, ext)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}
	}
	return fmt.Sprintf("%s-%s-99%s", base, t.Format("2006-01-02"), ext)
}

func (rl *RotatingLogger) compressOldLogs() error {
	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var logFiles []string
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, baseName) && strings.HasSuffix(name, ".log") {
			logFiles = append(logFiles, filepath.Join(dir, name))
		}
	}

	if len(logFiles) <= maxBackups {
		return nil
	}

	for i := 0; i < len(logFiles)-maxBackups; i++ {
		if compressOld {
			if err := compressFile(logFiles[i]); err != nil {
				log.Printf("Failed to compress %s: %v", logFiles[i], err)
			}
		} else {
			if err := os.Remove(logFiles[i]); err != nil {
				log.Printf("Failed to remove %s: %v", logFiles[i], err)
			}
		}
	}

	return nil
}

func compressFile(path string) error {
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

	// Simple compression simulation
	// In real implementation, use gzip.NewWriter(dst)
	_, err = io.Copy(dst, src)
	if err != nil {
		os.Remove(dstPath)
		return err
	}

	return os.Remove(path)
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
	logger, err := NewRotatingLogger("/var/log/myapp/application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "APP: ", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d with some sample data", i)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type LogRotator struct {
	mu          sync.Mutex
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
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
		filePath:    filePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		file:        file,
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

	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	ext := filepath.Ext(lr.filePath)
	base := lr.filePath[:len(lr.filePath)-len(ext)]
	timestamp := fmt.Sprintf("%d", os.Getpid())

	archivePath := fmt.Sprintf("%s-%s%s.gz", base, timestamp, ext)
	if err := compressFile(lr.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(lr.filePath); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	lr.currentSize = 0
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

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator("app.log", 10)
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: Application event recorded\n", i)
		if _, err := rotator.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
	}
}
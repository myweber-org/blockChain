
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

type RotatingLogger struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		backupCount: backupCount,
		compressOld: compressOld,
	}

	if err := rl.openOrCreate(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openOrCreate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
		rl.file = nil
	}

	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		
		if i == 0 {
			oldPath = rl.basePath
		} else {
			oldPath = filepath.Join(dir, fmt.Sprintf("%s.%d", baseName, i))
			if rl.compressOld {
				oldPath += ".gz"
			}
		}

		newPath = filepath.Join(dir, fmt.Sprintf("%s.%d", baseName, i+1))
		if rl.compressOld && i+1 < rl.backupCount {
			newPath += ".gz"
		}

		if _, err := os.Stat(oldPath); err == nil {
			if i+1 >= rl.backupCount {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if rl.compressOld {
		go rl.compressFile(rl.basePath, rl.basePath+".1.gz")
	} else {
		os.Rename(rl.basePath, rl.basePath+".1")
	}

	return rl.openOrCreate()
}

func (rl *RotatingLogger) compressFile(src, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	compressor := NewGzipCompressor(dstFile)
	_, err = io.Copy(compressor, srcFile)
	compressor.Close()

	if err == nil {
		os.Remove(src)
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

type GzipCompressor struct {
	w io.WriteCloser
}

func NewGzipCompressor(w io.Writer) *GzipCompressor {
	return &GzipCompressor{}
}

func (g *GzipCompressor) Write(p []byte) (int, error) {
	return len(p), nil
}

func (g *GzipCompressor) Close() error {
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: Application event at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
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
	dateStr := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != dateStr {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLogs(); err != nil {
				return err
			}
			if err := rl.cleanupOldBackups(); err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s.%s.log", rl.basePath, dateStr)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
		rl.currentDay = dateStr
	}
	return nil
}

func (rl *RotatingLogger) compressOldLogs() error {
	dir := filepath.Dir(rl.basePath)
	base := filepath.Base(rl.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		matched, err := filepath.Match(base+".*.log", name)
		if err != nil || !matched {
			continue
		}

		if filepath.Ext(name) == ".gz" {
			continue
		}

		srcPath := filepath.Join(dir, name)
		dstPath := srcPath + ".gz"

		if err := compressFile(srcPath, dstPath); err != nil {
			return err
		}

		os.Remove(srcPath)
	}
	return nil
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

func (rl *RotatingLogger) cleanupOldBackups() error {
	dir := filepath.Dir(rl.basePath)
	base := filepath.Base(rl.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var backups []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		matched, err := filepath.Match(base+".*.log.gz", name)
		if err != nil || !matched {
			continue
		}
		backups = append(backups, filepath.Join(dir, name))
	}

	if len(backups) <= maxBackups {
		return nil
	}

	for i := 0; i < len(backups)-maxBackups; i++ {
		os.Remove(backups[i])
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
	logger, err := NewRotatingLogger("./app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))))
		time.Sleep(10 * time.Millisecond)
	}
}
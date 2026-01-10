package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxSize    = 1024 * 1024 // 1 MB
	logDir     = "./logs"
	baseLog    = "app.log"
	timeFormat = "20060102-150405"
)

func rotateLog() error {
	logPath := filepath.Join(logDir, baseLog)
	info, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() < maxSize {
		return nil
	}

	timestamp := time.Now().Format(timeFormat)
	archiveName := fmt.Sprintf("%s.%s.gz", baseLog, timestamp)
	archivePath := filepath.Join(logDir, archiveName)

	src, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, src); err != nil {
		return err
	}

	if err := os.Truncate(logPath, 0); err != nil {
		return err
	}

	log.Printf("Rotated log to %s", archiveName)
	return nil
}

func ensureLogDir() error {
	return os.MkdirAll(logDir, 0755)
}

func main() {
	if err := ensureLogDir(); err != nil {
		log.Fatal(err)
	}

	logFile := filepath.Join(logDir, baseLog)
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	for i := 0; i < 5; i++ {
		if err := rotateLog(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
		log.Printf("Test log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(500 * time.Millisecond)
	}
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
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.size += int64(n)

	if rl.size >= maxFileSize {
		if err := rl.performRotation(); err != nil {
			return n, err
		}
	}
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")
	if rl.currentDay != today || rl.file == nil {
		if rl.file != nil {
			rl.file.Close()
		}
		rl.currentDay = today
		return rl.openCurrentFile()
	}
	return nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s-%s.log", rl.basePath, rl.currentDay)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	return nil
}

func (rl *RotatingLogger) performRotation() error {
	if rl.file != nil {
		rl.file.Close()
	}

	oldPath := fmt.Sprintf("%s-%s.log", rl.basePath, rl.currentDay)
	for i := backupCount - 1; i >= 0; i-- {
		var source, dest string
		if i == 0 {
			source = oldPath
		} else {
			source = fmt.Sprintf("%s-%s.log.%d.gz", rl.basePath, rl.currentDay, i-1)
		}
		dest = fmt.Sprintf("%s-%s.log.%d.gz", rl.basePath, rl.currentDay, i)

		if _, err := os.Stat(source); err == nil {
			if i == backupCount-1 {
				os.Remove(dest)
			} else {
				os.Rename(source, dest)
			}
		}
	}

	if err := compressFile(oldPath, fmt.Sprintf("%s-%s.log.0.gz", rl.basePath, rl.currentDay)); err != nil {
		return err
	}

	os.Remove(oldPath)
	rl.size = 0
	return rl.openCurrentFile()
}

func compressFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
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

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/application")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 100; i++ {
		log.Printf("Log entry number %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}
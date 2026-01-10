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
package main

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
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName    = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    basePath    string
    fileSize    int64
}

func NewLogRotator(logDir string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    logPath := filepath.Join(logDir, logFileName)
    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to get file info: %w", err)
    }

    return &LogRotator{
        currentFile: file,
        basePath:    logDir,
        fileSize:    info.Size(),
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.fileSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, fmt.Errorf("failed to rotate log: %w", err)
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.fileSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return fmt.Errorf("failed to close current log file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(lr.basePath, fmt.Sprintf("app.%s.log", timestamp))

    if err := os.Rename(filepath.Join(lr.basePath, logFileName), backupPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    lr.currentFile = file
    lr.fileSize = 0

    go lr.cleanupOldFiles()

    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    files, err := filepath.Glob(filepath.Join(lr.basePath, "app.*.log"))
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i, file := range files {
        if i >= maxBackupFiles {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) parseBackupTimestamp(filename string) (time.Time, error) {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) != 3 || parts[0] != "app" || parts[2] != "log" {
        return time.Time{}, fmt.Errorf("invalid backup filename format")
    }

    timestamp := parts[1]
    if len(timestamp) != 15 {
        return time.Time{}, fmt.Errorf("invalid timestamp length")
    }

    year, _ := strconv.Atoi(timestamp[0:4])
    month, _ := strconv.Atoi(timestamp[4:6])
    day, _ := strconv.Atoi(timestamp[6:8])
    hour, _ := strconv.Atoi(timestamp[9:11])
    minute, _ := strconv.Atoi(timestamp[11:13])
    second, _ := strconv.Atoi(timestamp[13:15])

    return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
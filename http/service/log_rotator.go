
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
	maxFileSize = 10 * 1024 * 1024
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	baseName  string
	fileSize  int64
	fileIndex int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.fileIndex = rl.findLatestIndex()
	filename := rl.generateFilename(rl.fileIndex)
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
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) findLatestIndex() int {
	pattern := filepath.Join(logDir, rl.baseName+"*.log")
	matches, _ := filepath.Glob(pattern)
	maxIndex := 0

	for _, match := range matches {
		var index int
		fmt.Sscanf(filepath.Base(match), rl.baseName+".%d.log", &index)
		if index > maxIndex {
			maxIndex = index
		}
	}

	return maxIndex
}

func (rl *RotatingLogger) generateFilename(index int) string {
	return filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, index))
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.fileSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
		rl.compressCurrentFile()
	}

	rl.fileIndex++
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressCurrentFile() {
	oldFile := rl.generateFilename(rl.fileIndex)
	compressedFile := oldFile + ".gz"

	go func() {
		src, err := os.Open(oldFile)
		if err != nil {
			return
		}
		defer src.Close()

		dst, err := os.Create(compressedFile)
		if err != nil {
			return
		}
		defer dst.Close()

		gz := gzip.NewWriter(dst)
		defer gz.Close()

		if _, err := io.Copy(gz, src); err != nil {
			return
		}

		os.Remove(oldFile)
	}()
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
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
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
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
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
        sequence:    0,
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
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.sequence)
    rl.sequence++

    if err := os.Rename(rl.basePath, archivedPath); err != nil {
        return err
    }

    if err := rl.compressFile(archivedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
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

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    os.Remove(source)
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

func findNextSequence(basePath string) int {
    pattern := filepath.Base(basePath) + ".*"
    matches, _ := filepath.Glob(pattern)

    maxSeq := 0
    for _, match := range matches {
        ext := filepath.Ext(match)
        if len(ext) > 1 {
            if seq, err := strconv.Atoi(ext[1:]); err == nil && seq > maxSeq {
                maxSeq = seq
            }
        }
    }
    return maxSeq + 1
}
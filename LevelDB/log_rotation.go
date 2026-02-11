package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    file        *os.File
}

func NewRotatingLogger(filePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filePath: filePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rl.openFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openFile() error {
    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) rotate() error {
    rl.file.Close()

    for i := rl.maxFiles - 1; i > 0; i-- {
        oldPath := rl.backupPath(i)
        newPath := rl.backupPath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    if err := os.Rename(rl.filePath, rl.backupPath(1)); err != nil {
        return err
    }

    return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
    if index == 0 {
        return rl.filePath
    }
    return rl.filePath + "." + strconv.Itoa(index)
}

func (rl *RotatingLogger) cleanup() error {
    dir := filepath.Dir(rl.filePath)
    base := filepath.Base(rl.filePath)

    files, err := filepath.Glob(filepath.Join(dir, base+".*"))
    if err != nil {
        return err
    }

    var backupFiles []string
    for _, file := range files {
        parts := strings.Split(file, ".")
        if len(parts) > 0 {
            if _, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
                backupFiles = append(backupFiles, file)
            }
        }
    }

    sort.Strings(backupFiles)

    if len(backupFiles) > rl.maxFiles {
        for _, file := range backupFiles[rl.maxFiles:] {
            os.Remove(file)
        }
    }

    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
        if err := rl.cleanup(); err != nil {
            fmt.Printf("Cleanup error: %v\n", err)
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
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
}package main

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
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
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
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	rl.currentNum++
	if rl.currentNum > maxBackups {
		rl.currentNum = 1
	}

	backupPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentNum)
	if err := compressFile(rl.basePath, backupPath); err != nil {
		return err
	}

	if err := os.Truncate(rl.basePath, 0); err != nil {
		return err
	}

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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation completed. Check app.log.*.gz files")
}
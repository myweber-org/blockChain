
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
    mu       sync.Mutex
    file     *os.File
    size     int64
    basePath string
    current  string
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
    now := time.Now().Format("2006-01-02")
    filename := fmt.Sprintf("%s-%s.log", rl.basePath, now)
    
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
    rl.size = info.Size()
    rl.current = filename
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
    
    backupName := fmt.Sprintf("%s.%d.gz", rl.current, time.Now().Unix())
    if err := compressFile(rl.current, backupName); err != nil {
        return err
    }
    
    os.Remove(rl.current)
    
    if err := cleanupOldBackups(rl.basePath); err != nil {
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

func cleanupOldBackups(basePath string) error {
    pattern := fmt.Sprintf("%s-*.log.*.gz", basePath)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, file := range toDelete {
            os.Remove(file)
        }
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
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
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

	path := rl.basePath + ".log"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	if err := rl.file.Close(); err != nil {
		return err
	}

	oldPath := rl.basePath + ".log"
	newPath := fmt.Sprintf("%s.%d.log.gz", rl.basePath, time.Now().Unix())

	if err := compressFile(oldPath, newPath); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	rl.cleanupOldBackups()
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

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := rl.basePath + ".*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		toRemove := matches[:len(matches)-maxBackups]
		for _, path := range toRemove {
			os.Remove(path)
		}
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

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}
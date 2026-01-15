
package main

import (
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
	if err := rl.file.Close(); err != nil {
		return err
	}

	for i := maxBackups - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(rl.basePath, rl.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) backupPath(num int) string {
	if num == 0 {
		return rl.basePath
	}
	ext := filepath.Ext(rl.basePath)
	base := rl.basePath[:len(rl.basePath)-len(ext)]
	return fmt.Sprintf("%s.%d%s", base, num, ext)
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
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
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
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    file        *os.File
    maxSize     int64
    basePath    string
    currentSize int64
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
        file:        file,
        maxSize:     maxSize,
        basePath:    basePath,
        currentSize: info.Size(),
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
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.%d.log", rl.basePath, timestamp, rl.sequence)
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
    reader, err := os.Open(source)
    if err != nil {
        return err
    }
    defer reader.Close()

    compressedPath := source + ".gz"
    writer, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer writer.Close()

    gzWriter := gzip.NewWriter(writer)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, reader); err != nil {
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
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
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
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	baseName  string
	sequence  int
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: filepath.Join(logDir, name),
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	filename := fmt.Sprintf("%s.log", rl.baseName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	rl.sequence++
	if rl.sequence > maxBackups {
		rl.sequence = 1
	}

	oldName := fmt.Sprintf("%s.log", rl.baseName)
	newName := fmt.Sprintf("%s.%d.log.gz", rl.baseName, rl.sequence)

	if err := rl.compressFile(oldName, newName); err != nil {
		return err
	}

	if err := os.Remove(oldName); err != nil {
		return err
	}

	return rl.openFile()
}

func (rl *RotatingLogger) compressFile(src, dst string) error {
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxFiles    = 5
    logDir      = "./logs"
)

type RotatingLogger struct {
    baseName   string
    current    *os.File
    fileNumber int
    fileSize   int64
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
    }

    if err := rl.openNextFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
    if rl.current != nil {
        rl.current.Close()
    }

    rl.fileNumber = (rl.fileNumber % maxFiles) + 1
    fileName := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileNumber))

    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }

    rl.current = file
    rl.fileSize = 0

    rl.cleanupOldFiles()
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
    if err != nil {
        return
    }

    if len(files) <= maxFiles {
        return
    }

    var fileInfos []struct {
        path string
        num  int
    }

    for _, file := range files {
        base := filepath.Base(file)
        suffix := strings.TrimPrefix(strings.TrimSuffix(base, ".log"), rl.baseName+"_")
        if num, err := strconv.Atoi(suffix); err == nil {
            fileInfos = append(fileInfos, struct {
                path string
                num  int
            }{path: file, num: num})
        }
    }

    for _, info := range fileInfos {
        if info.num != rl.fileNumber {
            os.Remove(info.path)
        }
    }
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.openNextFile(); err != nil {
            return 0, err
        }
    }

    n, err = rl.current.Write(p)
    rl.fileSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) Close() error {
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    for i := 0; i < 100; i++ {
        log.Printf("Log entry number %d at %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(50 * time.Millisecond)
    }

    fmt.Println("Log rotation completed. Check ./logs directory.")
}
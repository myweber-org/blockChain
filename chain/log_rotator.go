
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    baseName    string
    sequence    int
}

func NewRotatingWriter(baseName string) (*RotatingWriter, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    w := &RotatingWriter{
        baseName: baseName,
        sequence: 0,
    }

    if err := w.openCurrentFile(); err != nil {
        return nil, err
    }

    return w, nil
}

func (w *RotatingWriter) openCurrentFile() error {
    filename := filepath.Join(logDir, fmt.Sprintf("%s.log", w.baseName))
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.currentFile = file
    w.currentSize = stat.Size()
    return nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.currentFile.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    w.sequence++
    if w.sequence > maxBackups {
        w.sequence = 1
    }

    oldName := filepath.Join(logDir, fmt.Sprintf("%s.log", w.baseName))
    newName := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", w.baseName, w.sequence))

    if err := os.Rename(oldName, newName); err != nil && !os.IsNotExist(err) {
        return err
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        os.Exit(1)
    }
    defer writer.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := writer.Write([]byte(msg)); err != nil {
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
    currentFile *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    fileCount   int
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
        fileCount:   0,
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
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, archivedPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.fileCount++

    if rl.fileCount > 10 {
        rl.cleanOldArchives()
    }

    return nil
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLogger) cleanOldArchives() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > 10 {
        filesToDelete := matches[:len(matches)-10]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
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

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    currentSize int64
    maxSize     int64
    file        *os.File
    sequence    int
}

func NewRotatingLog(basePath string, maxSizeMB int64) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSizeMB * 1024 * 1024,
        sequence: 0,
    }

    if err := rl.openCurrent(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) openCurrent() error {
    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.file = file
    rl.currentSize = stat.Size()

    rl.findLatestSequence()
    return nil
}

func (rl *RotatingLog) findLatestSequence() {
    pattern := rl.basePath + ".*.gz"
    matches, _ := filepath.Glob(pattern)
    maxSeq := 0

    for _, match := range matches {
        base := filepath.Base(match)
        if len(base) > len(rl.basePath)+5 {
            seqStr := base[len(rl.basePath)+1 : len(base)-3]
            if seq, err := strconv.Atoi(seqStr); err == nil && seq > maxSeq {
                maxSeq = seq
            }
        }
    }

    rl.sequence = maxSeq
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
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

func (rl *RotatingLog) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedName := fmt.Sprintf("%s.%d.%s", rl.basePath, rl.sequence, timestamp)
    rl.sequence++

    if err := os.Rename(rl.basePath, rotatedName); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedName); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLog) compressFile(source string) error {
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

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("application.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
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
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
	rotateChan chan struct{}
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filename:   filename,
		rotateChan: make(chan struct{}, 1),
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	go rl.monitorRotation()

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		rl.current.Close()
	}

	f, err := os.OpenFile(rl.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.current = f
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err = rl.current.Write(p)
	if err != nil {
		return n, err
	}

	rl.size += int64(n)

	if rl.size >= maxFileSize {
		select {
		case rl.rotateChan <- struct{}{}:
		default:
		}
	}

	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current == nil {
		return fmt.Errorf("no open file")
	}

	if err := rl.current.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", rl.filename, timestamp)

	if err := os.Rename(rl.filename, backupName); err != nil {
		return err
	}

	if err := rl.compressFile(backupName); err != nil {
		log.Printf("Failed to compress %s: %v", backupName, err)
	}

	if err := rl.cleanupOldBackups(); err != nil {
		log.Printf("Failed to cleanup old backups: %v", err)
	}

	return rl.openFile()
}

func (rl *RotatingLogger) compressFile(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	pattern := rl.filename + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= backupCount {
		return nil
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
		}
	}

	return nil
}

func (rl *RotatingLogger) monitorRotation() {
	for range rl.rotateChan {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	close(rl.rotateChan)
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			log.Printf("Write failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
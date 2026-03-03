
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

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    basePath   string
    currentSize int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        file:       file,
        basePath:   path,
        currentSize: info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
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
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.cleanupOldBackups()

    return nil
}

func (rl *RotatingLogger) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
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

func (rl *RotatingLogger) cleanupOldBackups() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            backups = append(backups, name)
        }
    }

    if len(backups) <= maxBackups {
        return
    }

    sortBackups(backups)

    for i := 0; i < len(backups)-maxBackups; i++ {
        os.Remove(filepath.Join(dir, backups[i]))
    }
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }
    timestamp := parts[len(parts)-2]
    ts, _ := strconv.ParseInt(strings.ReplaceAll(timestamp, "_", ""), 10, 64)
    return ts
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("logs/app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
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
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	if info, err := file.Stat(); err == nil {
		rl.currentSize = info.Size()
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	oldLogs, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return err
	}

	if len(oldLogs) >= maxBackups {
		oldest := oldLogs[0]
		if err := compressAndRemove(oldest); err != nil {
			return err
		}
	}

	return rl.openCurrentFile()
}

func compressAndRemove(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	destName := filename + ".gz"
	dest, err := os.Create(destName)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(filename)
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

type RotatingLog struct {
    filename   string
    current   *os.File
    size      int64
    mu        sync.Mutex
}

func NewRotatingLog(filename string) (*RotatingLog, error) {
    rl := &RotatingLog{filename: filename}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) openCurrent() error {
    file, err := os.OpenFile(rl.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.current = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLog) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", rl.filename, timestamp)

    if err := os.Rename(rl.filename, backupName); err != nil {
        return err
    }

    if err := rl.compressBackup(backupName); err != nil {
        return err
    }

    if err := rl.cleanOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLog) compressBackup(filename string) error {
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

    os.Remove(filename)
    return nil
}

func (rl *RotatingLog) cleanOldBackups() error {
    pattern := rl.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    toDelete := matches[:len(matches)-maxBackups]
    for _, file := range toDelete {
        os.Remove(file)
    }

    return nil
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLog("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
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
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compress bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compress,
    }
    
    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (r *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(r.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    
    if err := os.Rename(r.basePath, backupPath); err != nil {
        return err
    }
    
    if r.compressOld {
        compressedPath := backupPath + ".gz"
        if err := compressFile(backupPath, compressedPath); err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }
    
    if err := r.cleanOldBackups(); err != nil {
        fmt.Printf("Warning: failed to clean old backups: %v\n", err)
    }
    
    return r.openCurrentFile()
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

func (r *LogRotator) cleanOldBackups() error {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)
    
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }
    
    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backups = append(backups, name)
        }
    }
    
    if len(backups) <= r.maxBackups {
        return nil
    }
    
    sortBackups(backups)
    
    for i := r.maxBackups; i < len(backups); i++ {
        path := filepath.Join(dir, backups[i])
        os.Remove(path)
    }
    
    return nil
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return 0
    }
    
    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }
    
    ts, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return ts
}

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n", 
            time.Now().Format(time.RFC3339), i)
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
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLog struct {
	mu          sync.Mutex
	file        *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	backupCount int
}

func NewRotatingLog(filePath string, maxSizeMB int, backupCount int) (*RotatingLog, error) {
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

	return &RotatingLog{
		file:        file,
		filePath:    filePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		backupCount: backupCount,
	}, nil
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

	dir := filepath.Dir(rl.filePath)
	base := filepath.Base(rl.filePath)
	timestamp := time.Now().Format("20060102_150405")

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		if i == 0 {
			oldPath = rl.filePath
			newPath = filepath.Join(dir, fmt.Sprintf("%s.%s.gz", base, timestamp))
		} else {
			oldPath = filepath.Join(dir, fmt.Sprintf("%s.%d.gz", base, i-1))
			newPath = filepath.Join(dir, fmt.Sprintf("%s.%d.gz", base, i))
		}

		if _, err := os.Stat(oldPath); err == nil {
			if i == 0 {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLog) compressFile(src, dst string) error {
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

func (rl *RotatingLog) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	log, err := NewRotatingLog("app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer log.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		log.Write([]byte(msg))
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
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
}

func NewLogRotator(basePath string, maxSize int64, maxBackups int) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath:   basePath,
        maxSize:    maxSize,
        maxBackups: maxBackups,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
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

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.cleanupOldBackups()

    return lr.openCurrentFile()
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

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= lr.maxBackups {
        return
    }

    sortByTimestamp := func(files []string) []string {
        var timestamps []string
        fileMap := make(map[string]string)

        for _, file := range files {
            parts := strings.Split(file, ".")
            if len(parts) < 3 {
                continue
            }
            ts := parts[len(parts)-2]
            if _, err := strconv.ParseInt(ts, 10, 64); err == nil {
                timestamps = append(timestamps, ts)
                fileMap[ts] = file
            }
        }

        for i := 0; i < len(timestamps); i++ {
            for j := i + 1; j < len(timestamps); j++ {
                if timestamps[i] > timestamps[j] {
                    timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
                }
            }
        }

        var sorted []string
        for _, ts := range timestamps {
            sorted = append(sorted, fileMap[ts])
        }
        return sorted
    }

    sortedFiles := sortByTimestamp(matches)
    filesToDelete := len(sortedFiles) - lr.maxBackups

    for i := 0; i < filesToDelete; i++ {
        os.Remove(sortedFiles[i])
    }
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
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
    rotator, err := NewLogRotator("app.log", 1024*1024, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed")
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
    mu            sync.Mutex
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
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)
    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    rl.rotationCount++
    if err := rl.compressOldLog(rotatedPath); err != nil {
        fmt.Printf("Failed to compress log: %v\n", err)
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    return nil
}

func (rl *RotatingLogger) compressOldLog(path string) error {
    if rl.rotationCount%3 != 0 {
        return nil
    }

    srcFile, err := os.Open(path)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destPath := path + ".gz"
    destFile, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    os.Remove(path)
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

func (rl *RotatingLogger) cleanupOldFiles(maxFiles int) {
    pattern := rl.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxFiles {
        filesToDelete := matches[:len(matches)-maxFiles]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    logger.cleanupOldFiles(5)
    fmt.Println("Log rotation completed. Check generated files.")
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
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationSeq  int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}
	
	return &LogRotator{
		file:        file,
		basePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		rotationSeq: 0,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("rotation failed: %w", err)
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
		return fmt.Errorf("failed to close current file: %w", err)
	}
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", lr.basePath, timestamp, lr.rotationSeq)
	
	if err := os.Rename(lr.basePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}
	
	file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	
	lr.file = file
	lr.currentSize = 0
	lr.rotationSeq++
	
	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write failed: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
}

func NewLogRotator(filePath string, maxSize int64, maxAge time.Duration) (*LogRotator, error) {
    rotator := &LogRotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxAge:   maxAge,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
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

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) rotateIfNeeded() error {
    if lr.currentSize >= lr.maxSize {
        return lr.performRotation()
    }

    info, err := os.Stat(lr.filePath)
    if err != nil {
        return err
    }

    if time.Since(info.ModTime()) > lr.maxAge {
        return lr.performRotation()
    }

    return nil
}

func (lr *LogRotator) performRotation() error {
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    if err := lr.openCurrentFile(); err != nil {
        return err
    }

    lr.cleanOldBackups()
    return nil
}

func (lr *LogRotator) cleanOldBackups() {
    cutoffTime := time.Now().Add(-lr.maxAge * 2)
    pattern := fmt.Sprintf("%s.*", lr.filePath)

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    for _, match := range matches {
        info, err := os.Stat(match)
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoffTime) {
            os.Remove(match)
        }
    }
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if err := lr.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(message))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation example completed")
}
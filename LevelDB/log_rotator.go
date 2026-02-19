package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackupAge = 30 * 24 * time.Hour
)

type RotatingLogger struct {
	currentFile *os.File
	filePath    string
	bytesWritten int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{filePath: path}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.bytesWritten+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	
	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.bytesWritten += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)
	
	if err := compressFile(rl.filePath, backupPath); err != nil {
		return err
	}
	
	if err := os.Remove(rl.filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	rl.bytesWritten = 0
	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer source.Close()
	
	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()
	
	gz := gzip.NewWriter(dest)
	defer gz.Close()
	
	_, err = io.Copy(gz, source)
	return err
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	
	rl.currentFile = file
	rl.bytesWritten = info.Size()
	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func cleanupOldBackups(logPath string) error {
	pattern := logPath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	
	cutoff := time.Now().Add(-maxBackupAge)
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(match)
		}
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
	mu          sync.Mutex
	currentFile *os.File
	currentSize int64
	basePath    string
	sequence    int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
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

	rl.currentFile = file
	rl.currentSize = info.Size()
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

	backupPath := fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	if err := rl.compressBackup(backupPath); err != nil {
		return err
	}

	rl.sequence = (rl.sequence + 1) % maxBackups
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressBackup(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	os.Remove(src)
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

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
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

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    filename   string
    current    *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filename: filename}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) openCurrent() error {
    info, err := os.Stat(rl.filename)
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    if info != nil {
        rl.size = info.Size()
    } else {
        rl.size = 0
    }

    file, err := os.OpenFile(rl.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    rl.current = file
    return nil
}

func (rl *RotatingLogger) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }

    timestamp := time.Now().Format("20060102-150405")
    backupName := fmt.Sprintf("%s.%s.gz", rl.filename, timestamp)

    if err := compressFile(rl.filename, backupName); err != nil {
        return err
    }

    if err := os.Remove(rl.filename); err != nil {
        return err
    }

    cleanupOldBackups(rl.filename)

    rl.size = 0
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

func cleanupOldBackups(baseName string) {
    pattern := fmt.Sprintf("%s.*.gz", baseName)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer func() {
        if logger.current != nil {
            logger.current.Close()
        }
    }()

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
    "strconv"
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    basePath      string
    maxSize       int64
    maxBackups    int
    compressOld   bool
    currentSize   int64
    currentFile   *os.File
    mu            sync.Mutex
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (r *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(r.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
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
        err := r.rotate()
        if err != nil {
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

    timestamp := time.Now().Unix()
    backupPath := fmt.Sprintf("%s.%d", r.basePath, timestamp)

    err := os.Rename(r.basePath, backupPath)
    if err != nil {
        return err
    }

    if r.compressOld {
        compressedPath := backupPath + ".gz"
        err = compressFile(backupPath, compressedPath)
        if err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }

    err = r.cleanupOldBackups()
    if err != nil {
        return err
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

func (r *LogRotator) cleanupOldBackups() error {
    pattern := r.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    var backupFiles []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") || isTimestampFile(match) {
            backupFiles = append(backupFiles, match)
        }
    }

    if len(backupFiles) <= r.maxBackups {
        return nil
    }

    filesToRemove := backupFiles[:len(backupFiles)-r.maxBackups]
    for _, file := range filesToRemove {
        os.Remove(file)
    }

    return nil
}

func isTimestampFile(path string) bool {
    parts := strings.Split(path, ".")
    if len(parts) < 2 {
        return false
    }

    timestampStr := parts[len(parts)-1]
    _, err := strconv.ParseInt(timestampStr, 10, 64)
    return err == nil
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
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
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

type RotatingLog struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    rotationCount int
}

func NewRotatingLog(filePath string, maxSizeMB int) (*RotatingLog, error) {
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
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
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

    timestamp := time.Now().Format("20060102_150405")
    archivePath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)

    if err := compressFile(rl.filePath, archivePath); err != nil {
        return err
    }

    if err := os.Remove(rl.filePath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.rotationCount++

    return nil
}

func compressFile(source, target string) error {
    sourceFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    targetFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    gzWriter := gzip.NewWriter(targetFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    return err
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    log, err := NewRotatingLog("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	backupCount = 5
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

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
				os.Remove(oldPath)
			}
		}
	}

	if err := os.Rename(rl.basePath, rl.getBackupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) getBackupPath(num int) string {
	if num == 0 {
		return rl.basePath
	}
	ext := filepath.Ext(rl.basePath)
	base := rl.basePath[:len(rl.basePath)-len(ext)]
	if num == 1 {
		return fmt.Sprintf("%s.1%s", base, ext)
	}
	return fmt.Sprintf("%s.%d.gz", base, num)
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
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
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
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    filename   string
    current   *os.File
    size      int64
    mu        sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filename: filename}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) openCurrent() error {
    info, err := os.Stat(rl.filename)
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    if err == nil {
        rl.size = info.Size()
    } else {
        rl.size = 0
    }

    file, err := os.OpenFile(rl.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    rl.current = file
    return nil
}

func (rl *RotatingLogger) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedFile := fmt.Sprintf("%s.%s", rl.filename, timestamp)
    if err := os.Rename(rl.filename, rotatedFile); err != nil {
        return err
    }

    compressedFile := rotatedFile + ".gz"
    if err := compressFile(rotatedFile, compressedFile); err != nil {
        return err
    }
    os.Remove(rotatedFile)

    if err := cleanupOldBackups(rl.filename); err != nil {
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

func cleanupOldBackups(baseFilename string) error {
    pattern := baseFilename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) > maxBackups {
        toRemove := matches[:len(matches)-maxBackups]
        for _, file := range toRemove {
            os.Remove(file)
        }
    }
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}
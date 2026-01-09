package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    FilePath    string
    MaxSize     int64
    MaxFiles    int
    RotateEvery time.Duration
    lastRotate  time.Time
}

func NewLogRotator(filePath string, maxSize int64, maxFiles int, rotateEvery time.Duration) *LogRotator {
    return &LogRotator{
        FilePath:    filePath,
        MaxSize:     maxSize,
        MaxFiles:    maxFiles,
        RotateEvery: rotateEvery,
        lastRotate:  time.Now(),
    }
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if err := lr.checkRotation(); err != nil {
        return 0, err
    }

    file, err := os.OpenFile(lr.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    return file.Write(p)
}

func (lr *LogRotator) checkRotation() error {
    now := time.Now()
    shouldRotate := false

    if lr.RotateEvery > 0 && now.Sub(lr.lastRotate) >= lr.RotateEvery {
        shouldRotate = true
        lr.lastRotate = now
    }

    if !shouldRotate && lr.MaxSize > 0 {
        if info, err := os.Stat(lr.FilePath); err == nil && info.Size() >= lr.MaxSize {
            shouldRotate = true
        }
    }

    if shouldRotate {
        return lr.performRotation()
    }
    return nil
}

func (lr *LogRotator) performRotation() error {
    for i := lr.MaxFiles - 1; i > 0; i-- {
        oldName := fmt.Sprintf("%s.%d", lr.FilePath, i)
        newName := fmt.Sprintf("%s.%d", lr.FilePath, i+1)

        if _, err := os.Stat(oldName); err == nil {
            if err := os.Rename(oldName, newName); err != nil {
                return err
            }
        }
    }

    if _, err := os.Stat(lr.FilePath); err == nil {
        backupName := fmt.Sprintf("%s.1", lr.FilePath)
        return os.Rename(lr.FilePath, backupName)
    }

    return nil
}

func (lr *LogRotator) Cleanup() error {
    for i := lr.MaxFiles + 1; ; i++ {
        fileName := fmt.Sprintf("%s.%d", lr.FilePath, i)
        if _, err := os.Stat(fileName); os.IsNotExist(err) {
            break
        }
        if err := os.Remove(fileName); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    rotator := NewLogRotator(
        "/var/log/app.log",
        10*1024*1024,
        5,
        time.Hour*24,
    )

    message := fmt.Sprintf("[%s] Application started\n", time.Now().Format(time.RFC3339))
    _, err := rotator.Write([]byte(message))
    if err != nil {
        fmt.Printf("Write error: %v\n", err)
        return
    }

    if err := rotator.Cleanup(); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }

    fmt.Println("Log rotation completed successfully")
}package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    baseName    string
    sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
        sequence: 0,
    }

    if err := rl.openNewFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    filename := filepath.Join(logDir, fmt.Sprintf("%s.log", rl.baseName))
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) rotateIfNeeded() error {
    if rl.currentSize < maxFileSize {
        return nil
    }

    oldPath := filepath.Join(logDir, fmt.Sprintf("%s.log", rl.baseName))
    newPath := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, rl.sequence))

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    rl.sequence++
    if rl.sequence > maxBackups {
        rl.cleanupOldFiles()
    }

    return rl.openNewFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
    for i := 0; i <= rl.sequence-maxBackups; i++ {
        oldFile := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, i))
        os.Remove(oldFile)
    }
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if err := rl.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err = rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
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

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
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

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
    gzipExt      = ".gz"
)

type RotatingLogger struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    currentSize int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: strings.TrimSuffix(basePath, logExtension),
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s_%s%s", rl.basePath, timestamp, logExtension)

    if err := os.Rename(rl.basePath+logExtension, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressOldLogs(); err != nil {
        return err
    }

    if err := rl.cleanupOldBackups(); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.basePath+logExtension, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) compressOldLogs() error {
    files, err := filepath.Glob(rl.basePath + "_*" + logExtension)
    if err != nil {
        return err
    }

    for _, file := range files {
        if strings.HasSuffix(file, gzipExt) {
            continue
        }

        gzFile := file + gzipExt
        if err := compressFile(file, gzFile); err != nil {
            return err
        }

        os.Remove(file)
    }

    return nil
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

func (rl *RotatingLogger) cleanupOldBackups() error {
    files, err := filepath.Glob(rl.basePath + "_*" + logExtension + gzipExt)
    if err != nil {
        return err
    }

    if len(files) <= maxBackups {
        return nil
    }

    sortableFiles := make([]string, len(files))
    copy(sortableFiles, files)

    for i := 0; i < len(sortableFiles)-maxBackups; i++ {
        os.Remove(sortableFiles[i])
    }

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
    logger, err := NewRotatingLogger("application.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
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
    if os.IsNotExist(err) {
        file, err := os.Create(rl.filename)
        if err != nil {
            return err
        }
        rl.current = file
        rl.size = 0
        return nil
    }
    if err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filename, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    rl.current = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    backupName := fmt.Sprintf("%s.%s.gz", rl.filename, timestamp)

    if err := compressFile(rl.filename, backupName); err != nil {
        return err
    }

    if err := cleanupOldBackups(rl.filename); err != nil {
        return err
    }

    os.Remove(rl.filename)
    return rl.openCurrent()
}

func compressFile(source, target string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    return err
}

func cleanupOldBackups(baseName string) error {
    pattern := baseName + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
    return nil
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
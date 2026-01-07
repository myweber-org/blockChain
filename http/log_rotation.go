package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 1024 * 1024 // 1MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentSize int64
    file        *os.File
    logger      *log.Logger
}

func NewRotatingLogger() (*RotatingLogger, error) {
    rl := &RotatingLogger{}
    if err := rl.openLogFile(); err != nil {
        return nil, err
    }
    rl.logger = log.New(rl.file, "", log.LstdFlags)
    return rl, nil
}

func (rl *RotatingLogger) openLogFile() error {
    info, err := os.Stat(logFileName)
    if err != nil && !os.IsNotExist(err) {
        return err
    }
    if info != nil {
        rl.currentSize = info.Size()
    }

    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    rl.file = file
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
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
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    if err := os.Rename(logFileName, backupName); err != nil {
        return err
    }

    if err := rl.openLogFile(); err != nil {
        return err
    }
    rl.currentSize = 0
    rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    files, err := filepath.Glob(logFileName + ".*")
    if err != nil {
        return
    }

    if len(files) <= maxBackupFiles {
        return
    }

    for i := 0; i < len(files)-maxBackupFiles; i++ {
        os.Remove(files[i])
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger()
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().String())))
        time.Sleep(10 * time.Millisecond)
    }
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
    maxFileSize   = 1024 * 1024 // 1MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentSize int64
    file        *os.File
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    fullPath := filepath.Join(path, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentSize: info.Size(),
        file:        file,
        basePath:    path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            log.Printf("Rotation failed: %v", err)
        }
    }

    n, err = rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.file.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(rl.basePath, logFileName)
    newPath := filepath.Join(rl.basePath, backupName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0

    go rl.cleanupOldFiles()

    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := filepath.Join(rl.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    var files []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }

        timestamp := parts[len(parts)-1]
        t, err := time.Parse("20060102_150405", timestamp)
        if err != nil {
            continue
        }

        files = append(files, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(files)-maxBackupFiles; i++ {
        oldestIdx := i
        for j := i + 1; j < len(files); j++ {
            if files[j].time.Before(files[oldestIdx].time) {
                oldestIdx = j
            }
        }
        os.Remove(files[oldestIdx].path)
        files[oldestIdx] = files[i]
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1024))
        time.Sleep(10 * time.Millisecond)
    }
}
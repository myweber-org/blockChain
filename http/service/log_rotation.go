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
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupLogs = 5
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
    if err == nil {
        rl.currentSize = info.Size()
    }

    rl.file, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    return err
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.file.Write(p)
    rl.currentSize += int64(n)
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
    rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    pattern := fmt.Sprintf("%s.*", logFileName)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackupLogs {
        toDelete := matches[:len(matches)-maxBackupLogs]
        for _, f := range toDelete {
            os.Remove(f)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger()
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logger.Write([]byte(fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
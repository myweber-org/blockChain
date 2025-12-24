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
}package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentSize  int64
	maxSize      int64
	basePath     string
	file         *os.File
	currentIndex int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
	}

	if err := rl.openOrCreateLog(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openOrCreateLog() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	path := rl.basePath
	if rl.currentIndex > 0 {
		path = rl.basePath + "." + string(rune('0'+rl.currentIndex))
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	rl.currentIndex++
	return rl.openOrCreateLog()
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
	logger, err := NewRotatingLogger("app.log", 1024*1024)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)
	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry number %d", i)
	}
}
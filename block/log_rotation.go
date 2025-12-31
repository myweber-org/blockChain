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
	maxBackupLogs = 5
	logFileName   = "app.log"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: basePath}
	if err := rl.openCurrentLog(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentLog() error {
	fullPath := filepath.Join(rl.basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.currentSize+int64(len(p)) > maxLogSize {
		if err := rl.rotateLog(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateLog() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(rl.basePath, logFileName)
	newPath := filepath.Join(rl.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.cleanupOldLogs(); err != nil {
		log.Printf("Failed to cleanup old logs: %v", err)
	}

	return rl.openCurrentLog()
}

func (rl *RotatingLogger) cleanupOldLogs() error {
	pattern := filepath.Join(rl.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackupLogs {
		return nil
	}

	oldestFirst := matches[:len(matches)-maxBackupLogs]
	for _, file := range oldestFirst {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger(".")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type RotatingWriter struct {
    currentFile *os.File
    filePath    string
    currentSize int64
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
    writer := &RotatingWriter{
        filePath: path,
    }
    
    if err := writer.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return writer, nil
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
    
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", w.filePath, timestamp)
    
    if err := os.Rename(w.filePath, backupPath); err != nil {
        return err
    }
    
    cleanupOldBackups(w.filePath)
    
    return w.openCurrentFile()
}

func (w *RotatingWriter) openCurrentFile() error {
    file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    w.currentFile = file
    w.currentSize = info.Size()
    return nil
}

func cleanupOldBackups(basePath string) {
    pattern := fmt.Sprintf("%s.*", basePath)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) > maxBackups {
        filesToRemove := matches[:len(matches)-maxBackups]
        for _, file := range filesToRemove {
            os.Remove(file)
        }
    }
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("application.log")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n", 
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}
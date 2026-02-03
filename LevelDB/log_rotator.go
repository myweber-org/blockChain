package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

func rotateLogs(logDir, archiveDir, prefix string, maxFiles int) error {
    files, err := os.ReadDir(logDir)
    if err != nil {
        return fmt.Errorf("failed to read log directory: %w", err)
    }

    var logFiles []string
    for _, file := range files {
        if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
            logFiles = append(logFiles, filepath.Join(logDir, file.Name()))
        }
    }

    if len(logFiles) <= maxFiles {
        return nil
    }

    if err := os.MkdirAll(archiveDir, 0755); err != nil {
        return fmt.Errorf("failed to create archive directory: %w", err)
    }

    for i := 0; i < len(logFiles)-maxFiles; i++ {
        oldPath := logFiles[i]
        timestamp := time.Now().Format("20060102_150405")
        newName := fmt.Sprintf("%s_%s", filepath.Base(oldPath), timestamp)
        newPath := filepath.Join(archiveDir, newName)

        srcFile, err := os.Open(oldPath)
        if err != nil {
            return fmt.Errorf("failed to open source file: %w", err)
        }

        dstFile, err := os.Create(newPath)
        if err != nil {
            srcFile.Close()
            return fmt.Errorf("failed to create destination file: %w", err)
        }

        if _, err := io.Copy(dstFile, srcFile); err != nil {
            srcFile.Close()
            dstFile.Close()
            return fmt.Errorf("failed to copy file content: %w", err)
        }

        srcFile.Close()
        dstFile.Close()

        if err := os.Remove(oldPath); err != nil {
            return fmt.Errorf("failed to remove old log file: %w", err)
        }

        fmt.Printf("Rotated %s to %s\n", oldPath, newPath)
    }

    return nil
}

func main() {
    logDir := "./logs"
    archiveDir := "./archive"
    prefix := "app_"
    maxFiles := 5

    if err := rotateLogs(logDir, archiveDir, prefix, maxFiles); err != nil {
        fmt.Printf("Error rotating logs: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Log rotation completed successfully")
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
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
		rl.compressOldFile()
	}

	filename := filepath.Join(logDir, fmt.Sprintf("%s_%03d.log", rl.baseName, rl.sequence))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	rl.sequence++

	return nil
}

func (rl *RotatingLogger) compressOldFile() {
	if rl.sequence <= 1 {
		return
	}

	oldFilename := filepath.Join(logDir, fmt.Sprintf("%s_%03d.log", rl.baseName, rl.sequence-2))
	newFilename := oldFilename + ".gz"

	oldFile, err := os.Open(oldFilename)
	if err != nil {
		return
	}
	defer oldFile.Close()

	newFile, err := os.Create(newFilename)
	if err != nil {
		return
	}
	defer newFile.Close()

	gzWriter := gzip.NewWriter(newFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return
	}

	os.Remove(oldFilename)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.openNewFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
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
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
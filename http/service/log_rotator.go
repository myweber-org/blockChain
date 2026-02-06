package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize   = 1024 * 1024 // 1MB
	logFileName  = "app.log"
	archiveDir   = "archives"
)

func rotateLogIfNeeded() error {
	info, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if info.Size() < maxLogSize {
		return nil
	}

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := filepath.Join(archiveDir, fmt.Sprintf("%s_%s", logFileName, timestamp))
	
	if err := os.Rename(logFileName, archiveName); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	newFile, err := os.Create(logFileName)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	newFile.Close()

	fmt.Printf("Log rotated: %s -> %s\n", logFileName, archiveName)
	return nil
}

func writeLog(message string) error {
	if err := rotateLogIfNeeded(); err != nil {
		return err
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}
	
	return nil
}

func main() {
	for i := 1; i <= 100; i++ {
		msg := fmt.Sprintf("Log entry number %d", i)
		if err := writeLog(msg); err != nil {
			fmt.Printf("Error writing log: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
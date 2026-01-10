package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxSize    = 1024 * 1024 // 1 MB
	logDir     = "./logs"
	baseLog    = "app.log"
	timeFormat = "20060102-150405"
)

func rotateLog() error {
	logPath := filepath.Join(logDir, baseLog)
	info, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() < maxSize {
		return nil
	}

	timestamp := time.Now().Format(timeFormat)
	archiveName := fmt.Sprintf("%s.%s.gz", baseLog, timestamp)
	archivePath := filepath.Join(logDir, archiveName)

	src, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, src); err != nil {
		return err
	}

	if err := os.Truncate(logPath, 0); err != nil {
		return err
	}

	log.Printf("Rotated log to %s", archiveName)
	return nil
}

func ensureLogDir() error {
	return os.MkdirAll(logDir, 0755)
}

func main() {
	if err := ensureLogDir(); err != nil {
		log.Fatal(err)
	}

	logFile := filepath.Join(logDir, baseLog)
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	for i := 0; i < 5; i++ {
		if err := rotateLog(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
		log.Printf("Test log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(500 * time.Millisecond)
	}
}
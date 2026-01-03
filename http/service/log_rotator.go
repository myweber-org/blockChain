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
	maxBackupFiles = 5
	logFileName   = "app.log"
)

type RotatingWriter struct {
	currentSize int64
	file        *os.File
}

func NewRotatingWriter() (*RotatingWriter, error) {
	w := &RotatingWriter{}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	info, err := os.Stat(logFileName)
	if err == nil {
		w.currentSize = info.Size()
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	if w.currentSize+int64(len(p)) > maxLogSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.currentSize += int64(n)
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	if err := os.Rename(logFileName, backupName); err != nil {
		return err
	}

	if err := w.openFile(); err != nil {
		return err
	}

	w.cleanupOldBackups()
	return nil
}

func (w *RotatingWriter) cleanupOldBackups() {
	pattern := fmt.Sprintf("%s.*", logFileName)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackupFiles {
		return
	}

	for i := 0; i < len(matches)-maxBackupFiles; i++ {
		os.Remove(matches[i])
	}
}

func (w *RotatingWriter) Close() error {
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter()
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, writer))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}
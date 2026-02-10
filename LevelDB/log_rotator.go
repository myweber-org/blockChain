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
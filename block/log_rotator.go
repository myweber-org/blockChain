
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
	backupCount = 5
)

type RotatingLogger struct {
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
	rotateChan chan struct{}
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filename:   filename,
		rotateChan: make(chan struct{}, 1),
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	go rl.monitorRotation()

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	file, err := os.OpenFile(rl.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.current = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err = rl.current.Write(p)
	if err != nil {
		return n, err
	}

	rl.size += int64(n)

	if rl.size >= maxFileSize {
		select {
		case rl.rotateChan <- struct{}{}:
		default:
		}
	}

	return n, nil
}

func (rl *RotatingLogger) monitorRotation() {
	for range rl.rotateChan {
		rl.rotate()
	}
}

func (rl *RotatingLogger) rotate() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size < maxFileSize {
		return
	}

	rl.current.Close()

	for i := backupCount - 1; i >= 0; i-- {
		oldName := rl.backupName(i)
		newName := rl.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if i == backupCount-1 {
				os.Remove(oldName)
			} else {
				os.Rename(oldName, newName)
			}
		}
	}

	os.Rename(rl.filename, rl.backupName(0))
	rl.compressBackup(0)
	rl.openFile()
}

func (rl *RotatingLogger) backupName(index int) string {
	if index == 0 {
		return rl.filename + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.filename, index)
}

func (rl *RotatingLogger) compressBackup(index int) error {
	srcName := rl.backupName(index)
	dstName := fmt.Sprintf("%s.%d.gz", rl.filename, index+1)

	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	_, err = io.Copy(gz, src)
	if err != nil {
		return err
	}

	os.Remove(srcName)
	return nil
}

func (rl *RotatingLogger) Close() error {
	close(rl.rotateChan)
	return rl.current.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
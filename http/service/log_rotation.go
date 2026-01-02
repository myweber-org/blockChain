
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 1024 * 1024 // 1MB
	maxBackups  = 5
)

type RotatingWriter struct {
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	w := &RotatingWriter{filename: filename}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) >= maxFileSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.current.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotatingWriter) openFile() error {
	file, err := os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.current = file
	w.size = stat.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.current != nil {
		w.current.Close()
	}

	for i := maxBackups - 1; i >= 0; i-- {
		oldName := w.backupName(i)
		newName := w.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if i+1 >= maxBackups {
				os.Remove(newName)
			} else {
				os.Rename(oldName, newName)
			}
		}
	}

	if err := os.Rename(w.filename, w.backupName(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return w.openFile()
}

func (w *RotatingWriter) backupName(i int) string {
	if i == 0 {
		return w.filename
	}
	return fmt.Sprintf("%s.%d", w.filename, i)
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.current != nil {
		return w.current.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n", 
			time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
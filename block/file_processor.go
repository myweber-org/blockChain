package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Data string
}

func worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d: %s\n", id, task.ID, task.Data)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	const numWorkers = 3
	const numTasks = 10

	taskChan := make(chan Task, numTasks)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChan, &wg)
	}

	for i := 1; i <= numTasks; i++ {
		taskChan <- Task{ID: i, Data: fmt.Sprintf("payload-%d", i)}
	}
	close(taskChan)

	wg.Wait()
	fmt.Println("All tasks completed")
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileJob struct {
	Path    string
	Content string
}

type Result struct {
	Path  string
	Lines int
	Err   error
}

func processFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return lineCount, err
	}

	return lineCount, nil
}

func worker(id int, jobs <-chan FileJob, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing: %s\n", id, job.Path)
		lines, err := processFile(job.Path)
		results <- Result{Path: job.Path, Lines: lines, Err: err}
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	filePatterns := []string{"*.txt", "*.log", "*.csv"}
	jobs := make(chan FileJob, 100)
	results := make(chan Result, 100)

	var wg sync.WaitGroup
	numWorkers := 4

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		for _, pattern := range filePatterns {
			matches, _ := filepath.Glob(pattern)
			for _, match := range matches {
				jobs <- FileJob{Path: match}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	totalLines := 0
	processedFiles := 0
	for result := range results {
		if result.Err != nil {
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Err)
			continue
		}
		fmt.Printf("File: %s, Lines: %d\n", result.Path, result.Lines)
		totalLines += result.Lines
		processedFiles++
	}

	fmt.Printf("\nSummary: Processed %d files, Total lines: %d\n", processedFiles, totalLines)
}package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  []string
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		fp.wg.Add(1)
		go func(ln int, text string) {
			defer fp.wg.Done()
			processed := fp.processLine(ln, text)
			fp.mu.Lock()
			fp.results = append(fp.results, processed)
			fp.mu.Unlock()
		}(lineNumber, line)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fp.wg.Wait()
	return nil
}

func (fp *FileProcessor) processLine(number int, text string) string {
	return fmt.Sprintf("Line %d: %d characters", number, len(text))
}

func (fp *FileProcessor) GetResults() []string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return append([]string{}, fp.results...)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	err := processor.ProcessFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	for _, result := range results {
		fmt.Println(result)
	}
}
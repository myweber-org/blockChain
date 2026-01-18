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
}
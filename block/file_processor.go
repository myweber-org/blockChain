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
}
package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID int
}

func (task *Task) Process() {
	fmt.Printf("Processed task: %d\n", task.ID)
	time.Sleep(5 * time.Second)
}

type WorkerPool struct {
	Tasks             []Task
	concurrentWorkers int
	tasksChan         chan Task
	wg                sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for task := range wp.tasksChan {
		task.Process()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run() {
	wp.tasksChan = make(chan Task, len(wp.Tasks))

	for _ = range wp.concurrentWorkers {
		go wp.worker()
	}

	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}
	close(wp.tasksChan)

	wp.wg.Wait()
}

func main() {
	tasks := make([]Task, 100)
	for i := range len(tasks) {
		tasks[i] = Task{ID: i}
	}

	wp := WorkerPool{
		Tasks:             tasks,
		concurrentWorkers: 5,
	}

	wp.Run()
	fmt.Println("All tasks processed")
}

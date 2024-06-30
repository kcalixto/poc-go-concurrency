package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

type Task struct {
	buffer        []int
	lastRowNumber int
}

var totalProcessed int

func httpRequest(task Task) {
	// Simulate HTTP request delay
	sleepTime := time.Duration(random(3, 5))*time.Second + time.Duration(random(1, 1000))*time.Millisecond
	time.Sleep(sleepTime)
	fmt.Printf("- Consumed %d in %s, total: %d\n", len(task.buffer), sleepTime.String(), task.lastRowNumber)
}

func random(min, max int) int {
	return min + rand.Intn(max-min)
}

func fillTasks(taskChan chan Task) {
	totalItemsToProcess := 100000
	maxItemsOnQueue := 200

	buffer := make([]int, 0, maxItemsOnQueue)
	for range totalItemsToProcess {
		for range maxItemsOnQueue {
			totalProcessed++
			buffer = append(buffer, totalProcessed)
			// Simulate delay in filling tasks
			time.Sleep(1 * time.Millisecond)
		}

		taskChan <- Task{
			buffer:        buffer,
			lastRowNumber: totalProcessed,
		}
		fmt.Println("+ New Batch")
	}
	close(taskChan)
}

func main() {
	// limit concurrences
	concurrentWorkers := 5
	askForNextWorker := make(chan struct{}, concurrentWorkers)

	var wg sync.WaitGroup
	taskChan := make(chan Task, concurrentWorkers*2)

	// Start goroutine for filling tasks
	go fillTasks(taskChan)

	for task := range taskChan {
		wg.Add(1)

		askForNextWorker <- struct{}{}
		go func() {
			defer wg.Done()
			httpRequest(task)
			<-askForNextWorker
		}()
	}

	wg.Wait()
	fmt.Println("All tasks processed")
}

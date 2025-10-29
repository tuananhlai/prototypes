package main

import (
	"flag"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	fair        = "fair"
	fairChannel = "fair-channel"
	unfair      = "unfair"
)

// go run main.go fair
// go run main.go unfair
func main() {
	flag.Parse()
	processingStyle := flag.Arg(0)

	numTasks := 100
	// A list of integers representing the time it takes to process a task.
	// The time it takes to process tasks[i] is i + 1 milliseconds.
	tasks := make([]int, 0, numTasks)

	for i := range numTasks {
		tasks = append(tasks, i+1)
	}

	numWorkers := 10

	switch processingStyle {
	case fair:
		processTasksFair(tasks, numWorkers)
	case fairChannel:
		processTasksFairChannel(tasks, numWorkers)
	case unfair:
		processTasksUnfair(tasks, numWorkers)
	default:
		log.Fatal("invalid processing style. usage: go run main.go fair|unfair")
	}
}

// Process the given list of tasks in an unfair manner. Each worker is assigned
// an equal number of tasks to process, regardless of how long the tasks takes to complete.
func processTasksUnfair(tasks []int, numWorkers int) {
	log.Println("Processing tasks unfairly...")
	var wg sync.WaitGroup

	numTasks := len(tasks)

	numTaskPerWorker := numTasks / numWorkers
	remainer := numTasks % numWorkers

	var currentIdx int
	for i := range numWorkers {
		wg.Add(1)
		size := numTaskPerWorker
		if i < remainer {
			size++
		}
		endIdx := currentIdx + size

		go func(threadID int, start int, end int) {
			startTime := time.Now()
			for i := start; i < end; i++ {
				processTask(tasks[i])
			}
			log.Println("Thread", threadID, "finished in", time.Since(startTime))
			wg.Done()
		}(i, currentIdx, endIdx)

		currentIdx = endIdx
	}

	wg.Wait()
}

// processTasksFair processes the given list of tasks in a fair manner. Each worker continues to process
// the latest unprocessed tasks until all tasks are finished.
func processTasksFair(tasks []int, numWorkers int) {
	log.Println("Processing tasks fairly using atomic counters...")

	var idx atomic.Int64
	var wg sync.WaitGroup

	numTasks := int64(len(tasks))

	for i := range numWorkers {
		wg.Add(1)
		go func(threadID int) {
			startTime := time.Now()

			for {
				// use a single atomic operation to add + retrieve current value to prevent race condition.
				currentIdx := idx.Add(1) - 1
				if currentIdx >= numTasks {
					break
				}
				task := tasks[currentIdx]
				processTask(task)
			}

			log.Println("Thread", threadID, "finished in", time.Since(startTime))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// processTasksFairChannel processes the given list of tasks in a fair manner. Each worker continues to process
// the latest unprocessed tasks until all tasks are finished. Make use of a Go channel instead of an atomic counter.
func processTasksFairChannel(tasks []int, numWorkers int) {
	log.Println("Processing tasks fairly using Go channels...")
	taskChan := make(chan int)

	var wg sync.WaitGroup
	for i := range numWorkers {
		wg.Add(1)
		go func(threadID int) {
			startTime := time.Now()

			for task := range taskChan {
				processTask(task)
			}

			log.Println("Thread", threadID, "finished in", time.Since(startTime))
			wg.Done()
		}(i)
	}

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	wg.Wait()
}

// Simulate a task that takes a given amount of time (ms) to complete.
func processTask(timeMs int) {
	time.Sleep(time.Millisecond * time.Duration(timeMs))
}

package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	numTasks := 100
	tasks := make([]int, 0, numTasks)

	for i := range numTasks {
		tasks = append(tasks, i)
	}

}

func processTasksUnfair(tasks []int) {
	var wg sync.WaitGroup

	numTasks := len(tasks)
	numWorker := 10

	numTaskPerWorker := numTasks / numWorker

	var startIdx, endIdx int
	for i := range numWorker {
		wg.Add(1)
		startIdx = i * numTaskPerWorker
		endIdx = min((i+1)*numTaskPerWorker, numTasks)

		go func(start int, end int) {
			startTime := time.Now()
			for i := start; i < end; i++ {
				processTask(tasks[i])
			}
			log.Println("Processed from", start, "to", end, "in", time.Since(startTime))
		}(startIdx, endIdx)
	}
}

func processTask(task int) {
	time.Sleep(time.Millisecond * time.Duration(task))
}

package main

import (
	"fmt"
	"sync"
	"time"
)

// Demonstration of a simple deadlock scenario using mutexes in Go.
func main() {
	var mu1 sync.Mutex
	var mu2 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		mu1.Lock()
		defer mu1.Unlock()

		// Gives thread 2 time to lock mu2.
		time.Sleep(1 * time.Second)

		mu2.Lock()
		defer mu2.Unlock()
		fmt.Println("Thread 1 finished.")
		wg.Done()
	}()

	go func() {
		mu2.Lock()
		defer mu2.Unlock()

		// Gives thread 1 time to lock mu1.
		time.Sleep(1 * time.Second)

		mu1.Lock()
		defer mu1.Unlock()
		fmt.Println("Thread 2 finished.")
		wg.Done()
	}()

	wg.Wait()
}

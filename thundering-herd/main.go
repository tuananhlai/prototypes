package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	cache := make(map[string]string)

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(threadID int) {
			if _, ok := cache["key"]; !ok {
				log.Println("Thread", threadID, "reading from database")
				cache["key"] = readFromDatabase("key")
			}

			log.Println("Thread", threadID, "read value:", cache["key"])

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func readFromDatabase(key string) string {
	time.Sleep(1 * time.Second)
	return key
}

package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	numThreads := 10

	repo := NewRepository()

	var wg sync.WaitGroup
	for i := range numThreads {
		wg.Add(1)
		go func(threadID int) {
			data := repo.GetData("key")
			log.Println("Thread", threadID, "read value:", data)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

type Repository struct {
	cache map[string]string
	mu    sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		cache: make(map[string]string),
	}
}

// GetData retrieves data from the cache or database.
func (r *Repository) GetData(key string) string {
	if _, ok := r.cache[key]; ok {
		return r.cache[key]
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Try reading from the cache again
	// to see if it has been updated by another thread.
	_, ok := r.cache[key]
	if !ok {
		r.cache[key] = r.getDataFromDatabase(key)
	}

	return r.cache[key]
}

// getDataFromDatabase simulates an expensive database read operation.
func (r *Repository) getDataFromDatabase(key string) string {
	log.Println("Reading from database")
	time.Sleep(1 * time.Second)
	return key
}

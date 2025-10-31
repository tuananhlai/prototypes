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
	cache            map[string]string
	lockManagerMutex sync.Mutex
	lockMap          map[string]*sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		cache:   make(map[string]string),
		lockMap: make(map[string]*sync.Mutex),
	}
}

// GetData retrieves data from the cache or database.
func (r *Repository) GetData(key string) string {
	if _, ok := r.cache[key]; ok {
		return r.cache[key]
	}

	mu := r.getLockForCacheKey(key)
	mu.Lock()
	defer mu.Unlock()

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

// getLockForCacheKey retrieves or creates a key-specific mutex
func (r *Repository) getLockForCacheKey(key string) *sync.Mutex {
	r.lockManagerMutex.Lock()
	defer r.lockManagerMutex.Unlock()

	if _, ok := r.lockMap[key]; !ok {
		r.lockMap[key] = &sync.Mutex{}
	}
	return r.lockMap[key]
}

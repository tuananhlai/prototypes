package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	connStr = "postgres://postgres:postgres@localhost:5432/prototype?sslmode=disable"
)

func benchmarkNoPool() {
	start := time.Now()
	var wg sync.WaitGroup

	for range 500 {
		wg.Go(func() {
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Fatalf("error opening new connection: %v", err)
				return
			}
			defer db.Close()

			_, err = db.Query("SELECT 1")
			if err != nil {
				log.Fatalf("error executing query: %v", err)
			}
		})
	}
	wg.Wait()
	fmt.Println(time.Since(start))
}

func main() {
	benchmarkNoPool()
}

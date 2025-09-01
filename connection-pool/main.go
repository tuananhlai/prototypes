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

func benchmarkWithPool() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening new connection: %v", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(10)

	var wg sync.WaitGroup
	for range 500 {
		wg.Go(func() {
			rows, err := db.Query("SELECT 1")
			if err != nil {
				log.Fatalf("error executing query: %v", err)
			}
			defer rows.Close()
		})
	}
	wg.Wait()
}

func benchmarkNoPool() {
	var wg sync.WaitGroup

	for range 500 {
		wg.Go(func() {
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Fatalf("error opening new connection: %v", err)
				return
			}
			defer db.Close()

			rows, err := db.Query("SELECT 1")
			if err != nil {
				log.Fatalf("error executing query: %v", err)
			}
			rows.Close()
		})
	}
	wg.Wait()
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	benchmarkWithPool()
	// benchmarkNoPool()
}

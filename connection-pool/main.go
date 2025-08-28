package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	connStr = "postgres://postgres:postgres@localhost:5432/prototype?sslmode=disable"
)

func benchmarkOneConnection() {
	start := time.Now()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)

	var wg sync.WaitGroup

	for range 2 {
		wg.Go(func() {
			_, err = db.Query("SELECT 1")
			if err != nil {
				fmt.Printf("error executing query: %v", err)
			}
		})
	}
	wg.Wait()
	fmt.Println(time.Since(start))
}

func benchmarkNoPool() {
	start := time.Now()
	var wg sync.WaitGroup

	for range 500 {
		wg.Go(func() {
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			_, err = db.Query("SELECT 1")
			if err != nil {
				panic(err)
			}
		})
	}
	wg.Wait()
	fmt.Println(time.Since(start))
}

func main() {
	benchmarkOneConnection()
}

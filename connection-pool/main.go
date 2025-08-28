package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func benchmarkNoPool() {
	connStr := "postgres://postgres:postgres@localhost:5432/prototype?sslmode=disable"

	for range 500 {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Query("SELECT 1")
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	benchmarkNoPool()
}

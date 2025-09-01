package main

import (
	"context"
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

type ConnectionPool struct {
	db       *sql.DB
	connChan chan *sql.Conn
}

func NewConnectionPool(dsn string) (*ConnectionPool, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	connChan := make(chan *sql.Conn, 10)

	ctx := context.Background()
	for range 10 {
		conn, err := db.Conn(ctx)
		if err != nil {
			return nil, err
		}
		connChan <- conn
	}

	return &ConnectionPool{
		db:       db,
		connChan: connChan,
	}, nil
}

func (c *ConnectionPool) Take() *sql.Conn {
	return <-c.connChan
}

func (c *ConnectionPool) Put(conn *sql.Conn) {
	c.connChan <- conn
}

func benchmarkWithPool() {
	pool, err := NewConnectionPool(connStr)
	if err != nil {
		log.Fatalf("error creating connection pool: %v", err)
	}

	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			conn := pool.Take()
			defer pool.Put(conn)
			rows, err := conn.QueryContext(context.Background(), "SELECT 1")
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

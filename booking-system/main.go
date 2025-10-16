package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	connStr = "postgres://postgres:postgres@localhost:5432/prototype?sslmode=disable"
)

// A minimal booking system to demonstrate the problem with concurrent database access.
// Summary:
// Naive approach: ~300ms, <10 seats assigned.
// Subquery approach: ~300ms, <10 seats assigned.
// Locked approach: ~1s, 100 seats assigned.
// Optimized locked approach: ~300ms, 100 seats assigned. However, the time somes time reached 1s for some reason.
func main() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	err = setupDatabase(db)
	if err != nil {
		log.Fatalf("error setting up database: %v", err)
	}

	numSeats := 100
	err = generateSeats(db, numSeats)
	if err != nil {
		log.Fatalf("error generating bookings: %v", err)
	}

	numCustomers := 100

	start := time.Now()

	var wg sync.WaitGroup
	for i := 1; i <= numCustomers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = bookSeatSubquery(db, i)
		}(i)
	}
	wg.Wait()

	fmt.Printf("Time taken to assign seats: %v\n", time.Since(start))

	bookedSeats, err := countBookedSeats(db)
	if err != nil {
		log.Fatalf("error counting booked seats: %v", err)
	}
	fmt.Printf("%d seats are assigned\n", bookedSeats)
}

// setupDatabase deletes and recreates the `bookings` table.
func setupDatabase(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS bookings")
	if err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE bookings (seat_id INT NOT NULL, customer_id INT DEFAULT NULL, PRIMARY KEY (seat_id))`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

// generateSeats inserts the given number of seats into the `bookings` table.
func generateSeats(db *sql.DB, numSeats int) error {
	values := make([]string, numSeats)
	for i := range numSeats {
		values[i] = fmt.Sprintf("(%d)", i+1)
	}
	query := fmt.Sprintf("INSERT INTO bookings (seat_id) VALUES %s", strings.Join(values, ","))

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("error inserting bookings: %v", err)
	}

	return nil
}

// bookSeatNaive assigns an random empty seat to the given customer without any locks.
func bookSeatNaive(db *sql.DB, customerID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	rows := tx.QueryRow("SELECT seat_id FROM bookings WHERE customer_id IS NULL LIMIT 1")

	var seatID int
	err = rows.Scan(&seatID)
	if err != nil {
		return fmt.Errorf("error scanning seat ID: %v", err)
	}

	_, err = tx.Exec("UPDATE bookings SET customer_id = $1 WHERE seat_id = $2", customerID, seatID)
	if err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}

	return tx.Commit()
}

// bookSeatLocked assigns an empty seat to the given customer. The selected empty seat is locked to prevent
// other customers from selecting it.
func bookSeatLocked(db *sql.DB, customerID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	rows := tx.QueryRow("SELECT seat_id FROM bookings WHERE customer_id IS NULL FOR UPDATE LIMIT 1")

	var seatID int
	err = rows.Scan(&seatID)
	if err != nil {
		return fmt.Errorf("error scanning seat ID: %v", err)
	}

	_, err = tx.Exec("UPDATE bookings SET customer_id = $1 WHERE seat_id = $2", customerID, seatID)
	if err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}

	return tx.Commit()
}

// bookSeatLockedOptimized uses a similar approach to bookSeatLocked, but with the SKIP LOCKED option to prevent
// unnecessary waiting for the same row.
// This method sometimes failed to assign 1 seat for some reason.
func bookSeatLockedOptimized(db *sql.DB, customerID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	rows := tx.QueryRow(`
		SELECT seat_id FROM bookings 
		WHERE customer_id IS NULL FOR UPDATE SKIP LOCKED LIMIT 1
	`)

	var seatID int
	err = rows.Scan(&seatID)
	if err != nil {
		return fmt.Errorf("error scanning seat ID: %v", err)
	}

	_, err = tx.Exec("UPDATE bookings SET customer_id = $1 WHERE seat_id = $2", customerID, seatID)
	if err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}

	return tx.Commit()
}

func bookSeatSubquery(db *sql.DB, customerID int) error {
	_, err := db.Exec("UPDATE bookings SET customer_id = $1 WHERE seat_id = (SELECT seat_id FROM bookings WHERE customer_id IS NULL LIMIT 1)", customerID)
	if err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}
	return nil
}

func countBookedSeats(db *sql.DB) (int, error) {
	rows := db.QueryRow("SELECT COUNT(*) FROM bookings WHERE customer_id IS NOT NULL")

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error scanning count: %v", err)
	}
	return count, nil
}

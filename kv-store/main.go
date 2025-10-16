package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	connStr = "postgres://postgres:postgres@localhost:5432/prototype?sslmode=disable"
)

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

	kvStore := NewKVStore(db)

	err = kvStore.Put("key1", "value1", time.Minute)
	if err != nil {
		log.Fatalf("error putting kv: %v", err)
	}

	value, err := kvStore.Get("key1")
	if err != nil {
		log.Fatalf("error getting kv: %v", err)
	}
	fmt.Println(value)

	value, err = kvStore.Get("key2")
	if err != nil {
		log.Printf("error getting key2: %v", err)
	}
	fmt.Println(value)

	err = kvStore.Del("key1")
	if err != nil {
		log.Fatalf("error deleting kv: %v", err)
	}

	value, err = kvStore.Get("key1")
	if err != nil {
		log.Printf("error getting key1: %v", err)
	}
	fmt.Println(value)
}

func setupDatabase(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS kv")
	if err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE kv (
		key VARCHAR(255) PRIMARY KEY, 
		value VARCHAR(255), 
		expires_at TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

type KVStore struct {
	db *sql.DB
}

func NewKVStore(db *sql.DB) *KVStore {
	return &KVStore{db: db}
}

func (k *KVStore) Put(key string, value string, expiration time.Duration) error {
	expiresAt := time.Now().Add(expiration)

	_, err := k.db.Exec(`
	INSERT INTO kv (key, value, expires_at) VALUES ($1, $2, $3) 
	ON CONFLICT (key) DO UPDATE SET value = $2, expires_at = $3
	`, key, value, expiresAt)
	if err != nil {
		return fmt.Errorf("error inserting kv: %v", err)
	}

	return nil
}

func (k *KVStore) Get(key string) (string, error) {
	var value string

	err := k.db.QueryRow(`
	SELECT value FROM kv WHERE key = $1 AND expires_at > NOW()
	`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("key not found")
		}
		return "", fmt.Errorf("error scanning value: %v", err)
	}
	return value, nil
}

func (k *KVStore) Del(key string) error {
	_, err := k.db.Exec(`
	UPDATE kv SET expires_at = NOW() WHERE key = $1
	`, key)
	if err != nil {
		return fmt.Errorf("error deleting kv: %v", err)
	}

	return nil
}

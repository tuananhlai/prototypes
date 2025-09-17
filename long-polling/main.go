package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	timeout = 10 * time.Second
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("/poll", pollHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("error initializing http server: %v", err)
	}
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	messageChan := make(chan string)
	go func() {
		randomDuration := time.Duration(rand.Intn(15)) * time.Second
		time.Sleep(randomDuration)

		messageChan <- fmt.Sprintf("Current time is: %s", time.Now())
	}()

	select {
	case msg := <-messageChan:
		w.Write([]byte(msg))
	case <-time.After(timeout):
		log.Println("Request timed out.")
		w.WriteHeader(http.StatusNoContent)
	case <-r.Context().Done():
		log.Println("Client disconnected.")
		return
	}
}

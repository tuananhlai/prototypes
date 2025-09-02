package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	port = ":8787"
)

func sse(w http.ResponseWriter, r *http.Request) {
	// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#sending_events_from_the_server
	w.Header().Set("Content-Type", "text/event-stream")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	log.Println("Client connected")

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		default:
			message := fmt.Sprintf("The server time is %v", time.Now().Format(time.RFC1123))
			fmt.Fprintf(w, "data: %s\n\n", message)

			flusher.Flush()

			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/sse", sse)

	log.Printf("HTTP server started on %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("error initializing http server: %v", err)
	}
}

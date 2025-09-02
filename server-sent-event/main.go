package main

import (
	"log"
	"net/http"
)

const (
	port = ":8787"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	log.Printf("HTTP server started on %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("error initializing http server: %v", err)
	}
}

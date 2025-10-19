package main

import (
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	handler := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
	mux.Handle("/assets/", handler)

	slog.Info("Starting http server", "port", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("error starting http server", "error", err)
	}
}

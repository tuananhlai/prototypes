package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/ws", wsHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection:", err)
	}
	defer conn.Close()
	log.Println("Client connected.")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected: ", err)
			return
		}

		log.Println("Received message:", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("error writing message:", err)
			return
		}
	}
}

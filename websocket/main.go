package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	chatHandler := NewChatHandler()
	mux.Handle("/ws", chatHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}

type ChatHandler struct {
	connections []*websocket.Conn
	mux         sync.RWMutex
	upgrader    websocket.Upgrader
}

func NewChatHandler() *ChatHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &ChatHandler{
		connections: []*websocket.Conn{},
		mux:         sync.RWMutex{},
		upgrader:    upgrader,
	}
}

func (c *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := c.initConn(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.closeConn(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected: ", err)
			return
		}
		if messageType != websocket.TextMessage {
			log.Printf("Invalid data received: %d\n", messageType)
			return
		}

		if err = c.broadcast(p); err != nil {
			log.Printf("error writing message\n: %v", err)
			return
		}
	}
}

// broadcast sends the given message to all active clients.
func (c *ChatHandler) broadcast(message []byte) error {
	c.mux.RLock()

	var err error
	for _, conn := range c.connections {
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			err = fmt.Errorf("error writing message to connection: %v", err)
			break
		}
	}

	c.mux.RUnlock()

	return err
}

// initConn initializes a websocket connection and put it into the connection list.
func (c *ChatHandler) initConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("error upgrading connection: %v", err)
	}

	c.addConn(conn)

	return conn, nil
}

func (c *ChatHandler) addConn(conn *websocket.Conn) {
	c.mux.Lock()
	c.connections = append(c.connections, conn)
	c.mux.Unlock()
}

func (c *ChatHandler) removeConn(conn *websocket.Conn) {
	c.mux.Lock()
	c.connections = slices.DeleteFunc(c.connections, func(c *websocket.Conn) bool {
		return c == conn
	})
	c.mux.Unlock()
}

// closeConn closes the given websocket connection and removes it from
// the list.
func (c *ChatHandler) closeConn(conn *websocket.Conn) error {
	c.removeConn(conn)
	return conn.Close()
}

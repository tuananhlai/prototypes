package main

import (
	"flag"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 8080, "port to listen on")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	log.Printf("listening on port %d", port)
}

type BroadcastController struct {
	service          *BroadcastService
	upgrader         websocket.Upgrader
	rdb              *redis.Client
	redisChannelName string
}

func (b *BroadcastController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := b.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b.service.AddConnection(conn)
	defer b.service.RemoveConnection(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading message:", err)
			break
		}
	}
}

type BroadcastService struct {
	connections []*websocket.Conn
	mux         sync.RWMutex
}

func NewBroadcastService() *BroadcastService {
	return &BroadcastService{
		connections: make([]*websocket.Conn, 0),
	}
}

func (s *BroadcastService) AddConnection(conn *websocket.Conn) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.connections = append(s.connections, conn)
}

func (s *BroadcastService) RemoveConnection(conn *websocket.Conn) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.connections = slices.DeleteFunc(s.connections, func(c *websocket.Conn) bool {
		return c == conn
	})

	return conn.Close()
}

func (s *BroadcastService) Broadcast(message []byte) error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for _, conn := range s.connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return err
		}
	}

	return nil
}

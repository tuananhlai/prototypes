package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	channelName    = "messages"
	subscriberRole = "subscriber"
	publisherRole  = "publisher"
)

// Running the prototype:
// 1. Start the subscriber: go run main.go subscriber
// 2. Publish a single message: go run main.go publisher
func main() {
	flag.Parse()
	role := flag.Arg(0)

	switch role {
	case subscriberRole:
		runSubscriber()
	case publisherRole:
		runPublisher()
	default:
		fmt.Println("Invalid role. Please use 'subscriber' or 'publisher'.")
	}
}

func runSubscriber() {
	globalCtx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pubsub := rdb.Subscribe(globalCtx, channelName)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Printf("Received message: %s\n", msg.Payload)
	}
}

func runPublisher() {
	globalCtx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := rdb.Publish(globalCtx, channelName, time.Now()).Err()
	if err != nil {
		log.Fatalf("failed to publish message: %v", err)
	}
	fmt.Println("Message published")
}

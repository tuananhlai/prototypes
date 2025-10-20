package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		slog.Error("error creating tcp connection", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close()

	// These two headers are critical for the HTTP request to be processed correctly.
	request := "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		slog.Error("error writing data to TCP connection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	rawRes, err := io.ReadAll(reader)
	if err != nil {
		slog.Error("error reading response", slog.String("error", err.Error()))
		os.Exit(1)
	}

	fmt.Println(string(rawRes))
}

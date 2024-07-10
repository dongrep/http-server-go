package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buff := make([]byte, 1024)

	byteData, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Could not read from connection")
	}

	data := string(buff[:byteData])

	fmt.Printf("Received from request: %s", data)

	requestParts := strings.Split(data, " ")

	path := requestParts[1]
	if path == "" {
		fmt.Println("Invalid path")
	}

	if path == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(path, "/echo/") {
		message := strings.Split(path, "/echo/")[1]
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(message)) + "\r\n\r\n" + message))
	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Could not write to connection")
		return
	}
}

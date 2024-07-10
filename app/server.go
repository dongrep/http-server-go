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

	pathParams := strings.Split(path, "/")
	if pathParams[1] == "" {
		fmt.Println("Invalid pathParams")
	}

	path = "/" + pathParams[1]

	switch path {
	case "/":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	case "/echo":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(pathParams[2])) + "\r\n\r\n" + pathParams[2]))

	case "/user-agent":
		requestFields := strings.Split(data, "\r\n")
		for _, field := range requestFields {
			if strings.Contains(field, "User-Agent") {
				fieldValue := strings.Split(field, ": ")
				_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(fieldValue[1])) + "\r\n\r\n" + fieldValue[1] + "\r\n"))
				break
			}
		}
	default:
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Could not write to connection")
	}
}

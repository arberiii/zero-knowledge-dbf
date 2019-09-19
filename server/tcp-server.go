package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var SetSize = uint(0)
var currentRandom = 0

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	myDbf := New(10, 10)
	for _, elem := range elements {
		myDbf.Add(elem, currentRandom)
	}
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	getBuffer := bytes.NewBuffer(buf)
	m, err := DecodeMessage(*getBuffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.BloomFilter.Equal(myDbf.b))
	currentRandom = randomNumber()
	fmt.Println(currentRandom)
	newM := Message{Nonce: currentRandom}
	b := newM.EncodeMessage()
	conn.Write(b)
	// Close the connection when you're done with it.
	conn.Close()
}

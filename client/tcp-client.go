package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	randomInt := 0
	for i := 0; i < 200; i++{
		// connect to this socket
		conn, _ := net.Dial("tcp", "127.0.0.1:3333")
		toWrite := currentDBF(randomInt)
		conn.Write(toWrite)
		// listen
		message := bufio.NewReader(conn)
		buff := make([]byte, 1024)
		_, errr := message.Read(buff)
		if errr != nil {
			fmt.Println("Error reading:", errr.Error())
		}
		getBuffer := bytes.NewBuffer(buff)
		m, err := DecodeMessage(*getBuffer)
		if err != nil {
			log.Panic(err)
		}
		randomInt = m.Nonce
	}
}

func currentDBF(random int) []byte {
	myDbf := New(10, 10)
	for _, elem := range elements {
		myDbf.Add(elem, random)
	}
	m := Message{BloomFilter: myDbf.b}
	b := m.EncodeMessage()
	return b
}

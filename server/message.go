package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/willf/bitset"
)

type Message struct {
	Nonce       int
	SetSize     uint
	BloomFilter *bitset.BitSet
}

// TODO: the encoding should write directly to the connection. NewEncoder takes a writer as parameter
func (m *Message) EncodeMessage() []byte {
	var messageByte bytes.Buffer
	enc := gob.NewEncoder(&messageByte)
	err := enc.Encode(m)
	if err != nil {
		panic(err)
	}
	return messageByte.Bytes()
}

func DecodeMessage(b bytes.Buffer) (*Message, error) {
	var m Message
	dec := gob.NewDecoder(&b)
	if err := dec.Decode(&m); err != nil {
		fmt.Println(string(b.Bytes()))
		return &Message{}, err
	}
	return &m, nil
}

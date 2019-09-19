package main

import (
	"crypto/sha512"
	"strconv"
)

// iHash returns the ith hashed value
func iHash(data int, i int) [sha512.Size256]byte {
	newData := append([]byte(strconv.Itoa(data))[:], byte(i))
	return sha512.Sum512_256(newData)
}

func hashElement(element []byte) [sha512.Size256]byte {
	return sha512.Sum512_256(element)
}

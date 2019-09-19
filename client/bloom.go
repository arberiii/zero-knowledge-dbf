package main

import (
	"crypto/sha512"
	"math"
	"math/rand"
	"time"

	"github.com/willf/bitset"
)

// system-wide false positive rate constant
const fpr = 0.1

// periodic task time in seconds, there the periodic task is sending the bloom filter to neighbor
const period = 3

// setSize is basically the length of elements, used an extra variable to understand it easier
var setSize = uint(0)
var elements = [][]byte{[]byte("test")}

//var logger *log.Logger
var finalSetSize int

type DistBF struct {
	b *bitset.BitSet
	m uint
	k uint
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 100000000
	return (rand.Intn(max-min+1) + min)
}

// New function return the DBF generated from the sizes of to peers
func New(n1, n2 uint) *DistBF {
	m, k := EstimateParameters(n1, n2)
	return &DistBF{m: m, k: k, b: bitset.New(m)}
}

// EstimateParameters estimates requirements for m and k.
// Based on https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
func EstimateParameters(n1, n2 uint) (m uint, k uint) {
	n := max(n1, n2)
	m = uint(math.Ceil(-1 * float64(n) * math.Log(fpr) / math.Pow(math.Log(2), 2)))
	k = uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

func max(x, y uint) uint {
	if x < y {
		return y
	}
	return x
}

func xor(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}

	c := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func xorHash(a, b [sha512.Size256]byte) [sha512.Size256]byte {
	var c [sha512.Size256]byte
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func compare(bc1, bc2 *bitset.BitSet) (bool, uint, uint) {
	//comparable := true
	firstBigger := bc1.Difference(bc2).Count()
	secondBigger := bc2.Difference(bc1).Count()
	if firstBigger != 0 && secondBigger != 0 {
		return false, firstBigger, secondBigger
	}
	return true, firstBigger, secondBigger
}

func (dbf *DistBF) hashOfXOR(data int) (ret [][sha512.Size256]byte) {
	for i := 0; i < int(dbf.k); i++ {
		ret = append(ret, iHash(data, i))
	}
	return
}

// addElementHash xor the result of function hashOfXOR with the hash of element (component wise)❤
func (dbf *DistBF) addElementHash(element []byte, hashes [][sha512.Size256]byte) {
	h := hashElement(element)
	for i := 0; i < len(hashes); i++ {
		hashes[i] = xorHash(hashes[i], h)
	}
}

// hashesModule find the location where to set bits in Bloom Filter
func (dbf *DistBF) hashesModulo(hashes [][sha512.Size256]byte) (ret []uint) {
	for _, hash := range hashes {
		ret = append(ret, byteModuloM(dbf.m, hash))
	}
	return
}

// Add element to DBF
// TODO: save hashes of elements
func (dbf *DistBF) Add(element []byte, random int) {
	hashes := dbf.hashOfXOR(random)
	dbf.addElementHash(element, hashes)
	locations := dbf.hashesModulo(hashes)
	for _, location := range locations {
		dbf.b.Set(location)
	}
}

func compareBytes(a, b []byte) bool {
	for i:=0; i < len(a); i++{
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
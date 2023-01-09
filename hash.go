package goestimators

import "github.com/spaolacci/murmur3"

// Hashes the given byte array into another, 128-bit (16-byte) byte array
func hash128(b []byte) []byte {
	h128 := murmur3.New128()
	return h128.Sum(b)
}

func hash64(b []byte) uint64 {
	return murmur3.Sum64(b)
}

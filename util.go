package goestimators

import "unsafe"

func bytesToUint64(b []byte) uint64 {
	return *((*uint64)(unsafe.Pointer(&b)))
}

func uint64ToBytes(x uint64, b []byte) {
	xb := (*[8]byte)(unsafe.Pointer(&x))
	copy(b, (*xb)[:])
}

func countBitsSetInBytes(buf []byte) uint {
	var count uint
	for _, b := range buf {
		var c uint
		for ; b != 0; c++ {
			b &= b - 1
		}
		count += c
	}
	return count
}

func countBitsSetInUint64(v uint64) uint {
	var c uint
	for ; v != 0; c++ {
		v &= v - 1
	}
	return c
}

func countRInUint64(x uint64) uint64 {
	return ^x & (x + 1)
}

func countTrailing1InUint64(x uint64) byte {
	return byte(countBitsSetInUint64(countRInUint64(x) - 1))
}

func countTrailing1InUint64Alt(x uint64) byte {
	var c uint
	for x != 0 {
		b := x & 1
		if b == 0 {
			break
		}
		c++
		x = x >> 1
	}
	return byte(c)
}

func isUintPowerOf2(x uint) bool {
	return (x & (x - 1)) == 0
}

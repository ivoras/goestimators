package goestimators

import "unsafe"

func bytesToUint64(b []byte) uint64 {
	return *((*uint64)(unsafe.Pointer(&b)))
}

func uint64ToBytes(x uint64, b []byte) {
	xb := (*[8]byte)(unsafe.Pointer(&x))
	copy(b, (*xb)[:])
}

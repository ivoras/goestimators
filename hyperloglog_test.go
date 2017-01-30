package gohyperloglog

import (
	"crypto/rand"
	"fmt"
	"testing"
	"unsafe"
)

func TestLlogLogSimple(t *testing.T) {
	ll, _ := NewLogLog(1024)
	fmt.Printf("len(ll.buckets)=%d, ll.bucmask=%x, ll.bucbits=%d\n", len(ll.buckets), ll.bucmask, ll.bucbits)
	var buf [8]byte
	for i := uint64(0); i < 10000; i++ {
		uint64ToBytes(i, buf[:])
		ll.Observe(buf[:])
	}
	est := ll.Estimate()
	if est >= 15000 || est <= 5000 {
		t.Errorf("The estimate is a bit off: %d", est)
	}
	fmt.Println("LogLogSimple_10k", est)
}

func TestLogLogRandom(t *testing.T) {
	ll, _ := NewLogLog(1024)
	var buf [8]byte
	const NumEntries = 100000
	for i := 0; i < NumEntries; i++ {
		n, err := rand.Read(buf[:])
		if err != nil {
			t.Errorf("rand.Read() returned an error: %v", err)
		}
		if n != len(buf) {
			t.Errorf("Couldn't read %d bytes from rand.Read(), read %d", len(buf), n)
		}
		ll.Observe(buf[:])
	}
	est := ll.Estimate()
	if est <= 50000 || est >= 150000 {
		t.Errorf("The estimate is a bit off: %d", est)
	}
	fmt.Println("LogLogRandom_100k", est)
}

func TestSuperLogLogRandom(t *testing.T) {
	ll, _ := NewLogLog(1024)
	var buf [8]byte
	const NumEntries = 1000000
	for i := 0; i < NumEntries; i++ {
		n, err := rand.Read(buf[:])
		if err != nil {
			t.Errorf("rand.Read() returned an error: %v", err)
		}
		if n != len(buf) {
			t.Errorf("Couldn't read %d bytes from rand.Read(), read %d", len(buf), n)
		}
		ll.Observe(buf[:])
	}
	est := ll.SuperEstimate()
	if est <= 500000 || est >= 1500000 {
		t.Errorf("The estimate is a bit off: %d", est)
	}
	fmt.Println("SuperLogLogRandom_1M", est)
}

func TestBitsSet(t *testing.T) {
	if countTrailing1InUint64(uint64(23)) != 3 { // 10111
		t.Errorf("Error for %d: %d", 23, countTrailing1InUint64(uint64(23)))
	}
	if countTrailing1InUint64(uint64(255)) != 8 { // 11111111
		t.Errorf("Error for %d: %d", 255, countTrailing1InUint64(uint64(255)))
	}
	if countTrailing1InUint64(uint64(0)) != 0 { // 0
		t.Errorf("Error for %d: %d", 0, countTrailing1InUint64(uint64(0)))
	}
	if countTrailing1InUint64(uint64(897)) != 1 { // 1110000001
		t.Errorf("Error for %d: %d", 897, countTrailing1InUint64(uint64(897)))
	}
}

func TestBitsSetAlt(t *testing.T) {
	if countTrailing1InUint64Alt(uint64(23)) != 3 { // 10111
		t.Errorf("Error for %d: %d", 23, countTrailing1InUint64Alt(uint64(23)))
	}
	if countTrailing1InUint64Alt(uint64(255)) != 8 { // 11111111
		t.Errorf("Error for %d: %d", 255, countTrailing1InUint64Alt(uint64(255)))
	}
	if countTrailing1InUint64Alt(uint64(0)) != 0 { // 0
		t.Errorf("Error for %d: %d", 0, countTrailing1InUint64Alt(uint64(0)))
	}
	if countTrailing1InUint64Alt(uint64(897)) != 1 { // 1110000001
		t.Errorf("Error for %d: %d", 897, countTrailing1InUint64Alt(uint64(897)))
	}
}

func TestBitsSetBoth(t *testing.T) {
	var buf [8]byte
	for i := 0; i < 100000; i++ {
		n, err := rand.Read(buf[:])
		if err != nil {
			t.Errorf("rand.Read() returned an error: %v", err)
		}
		if n != len(buf) {
			t.Errorf("Couldn't read %d bytes from rand.Read(), read %d", len(buf), n)
		}
		x := bytesToUint64(buf[:])
		if countTrailing1InUint64(x) != countTrailing1InUint64Alt(x) {
			t.Errorf("Error at %v: %d vs %d", x, countTrailing1InUint64(x), countTrailing1InUint64Alt(x))
		}
	}

}

func BenchmarkCountTrailing1(b *testing.B) {
	var buf [8]byte
	for i := 0; i < b.N; i++ {
		n, err := rand.Read(buf[:])
		if err != nil {
			b.Errorf("rand.Read() returned an error: %v", err)
		}
		if n != len(buf) {
			b.Errorf("Couldn't read %d bytes from rand.Read(), read %d", len(buf), n)
		}
		x := bytesToUint64(buf[:])
		countTrailing1InUint64(x)
	}
}

func BenchmarkCountTrailing1Slow(b *testing.B) {
	var buf [8]byte
	for i := 0; i < b.N; i++ {
		n, err := rand.Read(buf[:])
		if err != nil {
			b.Errorf("rand.Read() returned an error: %v", err)
		}
		if n != len(buf) {
			b.Errorf("Couldn't read %d bytes from rand.Read(), read %d", len(buf), n)
		}
		x := bytesToUint64(buf[:])
		countTrailing1InUint64Alt(x)
	}
}

func bytesToUint64(b []byte) uint64 {
	return *((*uint64)(unsafe.Pointer(&b)))
}

func uint64ToBytes(x uint64, b []byte) {
	p := uintptr(unsafe.Pointer(&x))
	for i := 0; i < 8; i++ {
		b[i] = byte(*(*byte)(unsafe.Pointer(p)))
		p++
	}
}

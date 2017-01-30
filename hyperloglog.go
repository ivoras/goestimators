package gohyperloglog

import (
	"fmt"
	"math"
	"sort"

	"github.com/spaolacci/murmur3"
)

type byteSlice []byte

func (a byteSlice) Len() int           { return len(a) }
func (a byteSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byteSlice) Less(i, j int) bool { return a[i] < a[j] }

// LogLog is the structure with book-keeping data for the LogLog calculations.
type LogLog struct {
	buckets []byte
	bucmask uint
	bucbits byte
}

// NewLogLog returns an initialised LogLog structure.
func NewLogLog(nBuckets uint) (*LogLog, error) {
	if !isUintPowerOf2(nBuckets) {
		return nil, fmt.Errorf("%d is not a power of 2", nBuckets)
	}
	var ll LogLog
	ll.buckets = make([]byte, nBuckets)
	ll.bucmask = nBuckets - 1
	ll.bucbits = countTrailing1InUint64Alt(uint64(ll.bucmask))
	return &ll, nil
}

// Observe makes the LogLog calculation for the given byte array, "counting" it
// in its data structure.
func (l *LogLog) Observe(b []byte) {
	h := hash64(b)
	bucket := h & uint64(l.bucmask)
	count1 := countTrailing1InUint64Alt(h >> l.bucbits)
	if count1 > l.buckets[bucket] {
		l.buckets[bucket] = count1
	}
}

// Estimate returns a LogLog estimation of the cardinality of the given inputs
func (l *LogLog) Estimate() uint64 {
	var sum uint
	for _, b := range l.buckets {
		sum += uint(b)
	}
	average := float64(sum) / float64(len(l.buckets))
	return uint64(math.Pow(2, average) * float64(len(l.buckets)) * 0.79402)
}

// SuperEstimate returns a cardinality estimation based on the SuperLogLog modification of the LogLog algorithm
func (l *LogLog) SuperEstimate() uint64 {
	var sortedBuckets []byte
	copy(sortedBuckets[:], l.buckets[:])
	sort.Sort(byteSlice(sortedBuckets[:]))
	cutoff := int(float64(len(l.buckets)) * 0.9)
	var sum uint
	for i, b := range l.buckets {
		if i > cutoff {
			break
		}
		sum += uint(b)
	}
	average := float64(sum) / float64(cutoff)
	return uint64(math.Pow(2, average) * float64(cutoff) * 0.79402)
}

// HyperLogLog is the structure with book-keeping data for the HyperLogLog calculations.
type HyperLogLog struct {
	x int
}

// Observe makes the LogLog calculation for the given byte array, "counting" it
// in its data structure.
func (hl *HyperLogLog) Observe(b []byte) {

}

// ----------------------------------------------------

// Hashes the given byte array into another, 128-bit (16-byte) byte array
func hash128(b []byte) []byte {
	h128 := murmur3.New128()
	return h128.Sum(b)
}

func hash64(b []byte) uint64 {
	return murmur3.Sum64(b)
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

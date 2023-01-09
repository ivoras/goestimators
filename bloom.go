package goestimators

import "log"

// Bloom size enum
type bloomSize int

const (
	// Size256 indicates a 256-bit Bloom filter
	Size256 bloomSize = iota
	// Size65536 indicates a 65536-bit Bloom filter
	Size65536
)

// Bloom is the bloom filter structure
type Bloom struct {
	size  bloomSize
	nInts int
	base  []uint64
}

// NewBloom creates a new Bloom filter structure
func NewBloom(size bloomSize) *Bloom {
	var nInts int
	switch size {
	case Size256:
		nInts = 4
	case Size65536:
		nInts = 1024
	default:
		log.Panic("Invalid size:", size)
	}
	b := Bloom{
		nInts: nInts,
		size:  size,
		base:  make([]uint64, nInts),
	}
	return &b
}

// Add adds a buffer / item to the Bloom filter
func (b *Bloom) Add(buf []byte) {
	ints := b.buf2ints(buf)
	for i, ii := range ints {
		b.base[i] |= ii
	}
}

// Check checks if the given buffer / item is in the Bloom filter. Returns false if there's no chance that the item is in the buffer.
func (b *Bloom) Check(buf []byte) bool {
	ints := b.buf2ints(buf)
	for i, ii := range ints {
		if b.base[i]&ii != ii {
			return false
		}
	}
	return true
}

func (b *Bloom) buf2ints(buf []byte) []uint64 {
	i := make([]uint64, b.nInts)
	h := hash64(buf)

	var mask uint64
	var shift uint
	
	// A performance optimization: We either look at the first byte of the 64-bit hash and
	// set the single bit in the 256-bit Bloom table, or we look at the first two bytes
	// and set a single bit in the 65536-bit Bloom table.
	if b.size == Size256 {
		mask = 0xff
		shift = 8
		h = h >> 56 // 64 bits are too much for a 256-bit Bloom base
	} else if b.size == Size65536 {
		mask = 0xffff
		shift = 16
		h = h >> 32
	} else {
		log.Panic("Invalid size:", b.size)
	}

	for h > 0 {
		b := h & mask
		//fmt.Println(buf, h, b, shift)
		// set the b'th bit in i
		i[b/64] |= 1 << (b % 64)
		h = h >> shift
	}

	return i
}

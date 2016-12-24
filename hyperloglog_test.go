package gohyperloglog

import "testing"

func TestLlogLogSimple(t *testing.T) {
	ll := NewLogLog()
	in1 := [...]byte{0, 0xff, 9}
	ll.Observe(in1[:])
	in2 := [...]byte{0, 0xff, 100}
	ll.Observe(in2[:])
	in3 := [...]byte{0, 0xff, 1}
	ll.Observe(in3[:])
	in4 := [...]byte{0, 0xff, 2}
	ll.Observe(in4[:])
	in5 := [...]byte{0, 0xff, 3}
	ll.Observe(in5[:])
	if ll.Estimate() != 5 {
		t.Errorf("Expecting 5 uniques, got %d", ll.Estimate())
	}
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

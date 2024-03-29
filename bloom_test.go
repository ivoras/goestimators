package goestimators

import (
	"fmt"
	"testing"
)

func TestBloom256(t *testing.T) {
	var N uint64 = 100
	b := NewBloom(Size256)

	var i uint64
	buf := make([]byte, 8)
	for i = 0; i < N; i++ {
		uint64ToBytes(i, buf)
		//fmt.Println(i, buf)
		b.Add(buf)
	}

	notfound := 0
	for i = 0; i < N+10; i++ {
		uint64ToBytes(i, buf)
		if !b.Check(buf) {
			notfound++
		}
	}

	fmt.Println("false.prob:", b.FalsePositiveProbability(N), "notfound.256:", notfound)
	if notfound > 0 {
		t.Log("Integers in range ", N, " not found:", notfound)
	}
}

func TestBloom65536(t *testing.T) {
	var N uint64 = 30000
	b := NewBloom(Size65536)

	var i uint64
	buf := make([]byte, 8)
	for i = 0; i < N; i++ {
		uint64ToBytes(i, buf)
		//fmt.Println(i, buf)
		b.Add(buf)
	}

	notfound := 0
	for i = 0; i < N+10; i++ {
		uint64ToBytes(i, buf)
		if !b.Check(buf) {
			notfound++
		}
	}

	fmt.Println("false.prob:", b.FalsePositiveProbability(N), "notfound.65536:", notfound)
	if notfound > 0 {
		t.Log("Integers in range ", N, " not found:", notfound)
	}
}

func BenchmarkBloomObserve256(b *testing.B) {
	b.StopTimer()
	bl := NewBloom(Size256)
	b.StartTimer()

	var i uint64
	buf := make([]byte, 8)
	for i = 0; i < uint64(b.N); i++ {
		uint64ToBytes(i, buf)
		bl.Observe(buf)
	}
}

func BenchmarkBloomCheck256(b *testing.B) {
	b.StopTimer()
	bl := NewBloom(Size256)

	var i uint64
	buf := make([]byte, 8)
	for i = 0; i < uint64(b.N); i++ {
		uint64ToBytes(i, buf)
		bl.Observe(buf)
	}

	b.StartTimer()
	notfound := 0
	for i = 0; i < uint64(b.N)+10; i++ {
		uint64ToBytes(i, buf)
		if !bl.Check(buf) {
			notfound++
		}
	}
}

func BenchmarkBloomObserve65536(b *testing.B) {
	b.StopTimer()
	bl := NewBloom(Size65536)
	b.StartTimer()

	var i uint64
	buf := make([]byte, 8)
	for i = 0; i < uint64(b.N); i++ {
		uint64ToBytes(i, buf)
		bl.Observe(buf)
	}
}

func BenchmarkBloomCheck65536(b *testing.B) {
	b.StopTimer()
	bl := NewBloom(Size65536)

	var i uint64
	buf := make([]byte, 8)
	for i = 0; i < uint64(b.N); i++ {
		uint64ToBytes(i, buf)
		bl.Observe(buf)
	}

	b.StartTimer()
	notfound := 0
	for i = 0; i < uint64(b.N)+10; i++ {
		uint64ToBytes(i, buf)
		if !bl.Check(buf) {
			notfound++
		}
	}
}

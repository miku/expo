package kv

import (
	"hash/fnv"
	"testing"
)

func BenchmarkMapAccess(b *testing.B) {
	m := make(map[string]string)
	m["k"] = "value"
	for i := 0; i < b.N; i++ {
		_ = m["k"]
	}
}

func BenchmarkSliceAccess(b *testing.B) {
	s := make([]string, 8)
	s[5] = "value"
	for i := 0; i < b.N; i++ {
		_ = s[5]
	}
}

func BenchmarkFNV(b *testing.B) {
	h := fnv.New32()
	data := []byte("k")
	for i := 0; i < b.N; i++ {
		_ = h.Sum(data)
	}
}

func BenchmarkQuasi1(b *testing.B) {
	s := make([]string, 8)
	s[customFunc1("k", 8)] = "value"
	for i := 0; i < b.N; i++ {
		_ = s[customFunc1("k", 8)]
	}
}

func BenchmarkQuasi2(b *testing.B) {
	s := make([]string, 8)
	s[customFunc2("k", 8)] = "value"
	for i := 0; i < b.N; i++ {
		_ = s[customFunc2("k", 8)]
	}
}

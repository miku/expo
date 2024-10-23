package mapkey

import (
	"encoding/binary"
	"testing"
)

func BenchmarkByteSliceToString(b *testing.B) {
	data := []byte("Berlin")
	for i := 0; i < b.N; i++ {
		ByteSliceToString(data)
	}
}

func BenchmarkByteSliceToFNV(b *testing.B) {
	data := []byte("Berlin")
	for i := 0; i < b.N; i++ {
		ByteSliceToFNV(data)
	}
}

func BenchmarkByteSliceToIntKey(b *testing.B) {
	data := []byte("Berlin")
	for i := 0; i < b.N; i++ {
		ByteSliceToIntKey(data)
	}
}

func BenchmarkByteSliceToIntBinary(b *testing.B) {
	data := []byte("Berlin")
	for i := 0; i < b.N; i++ {
		ByteSliceToBigInt(data)
	}
}

func BenchmarkByteSliceToStringZeroAlloc(b *testing.B) {
	data := []byte("Berlin")
	for i := 0; i < b.N; i++ {
		ByteSliceToStringZeroAlloc(data)
	}
}

func BenchmarkStringToIndex16(b *testing.B) {
	data := "Berlin"
	for i := 0; i < b.N; i++ {
		StringToIndex16(data)
	}
}

func BenchmarkDjb2(b *testing.B) {
	data := "Berlin"
	for i := 0; i < b.N; i++ {
		Djb2(data)
	}
}

func BenchmarkFxHasher(b *testing.B) {
	data := "BerlinXX" // padded
	word := binary.NativeEndian.Uint64([]byte(data))
	for i := 0; i < b.N; i++ {
		FxHasher(word)
	}
}

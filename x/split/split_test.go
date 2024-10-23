package split

import (
	"fmt"
	"testing"
)

func BenchmarkStringSplit(b *testing.B) {
	var bms = []struct {
		s string
	}{
		{"Nuuk;14.3"},
		{"Las Palmas de Gran Canaria;14.3"},
	}
	for _, bm := range bms {
		name := fmt.Sprintf("s-len-%d", len(bm.s))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = StringSplit(bm.s)
			}
		})
	}
}

func BenchmarkByteSliceSplit(b *testing.B) {
	var bms = []struct {
		b []byte
	}{
		{[]byte("Nuuk;14.3")},
		{[]byte("Las Palmas de Gran Canaria;14.3")},
	}
	for _, bm := range bms {
		name := fmt.Sprintf("b-len-%d", len(bm.b))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = BytesSplit(bm.b)
			}
		})
	}
}

func BenchmarkCustomSplit(b *testing.B) {
	var bms = []struct {
		b []byte
	}{
		{[]byte("Nuuk;14.3")},
		{[]byte("Las Palmas de Gran Canaria;14.3")},
	}
	for _, bm := range bms {
		name := fmt.Sprintf("b-len-%d", len(bm.b))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = CustomSplit(bm.b)
			}
		})
	}
}

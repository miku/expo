package findchar

import (
	"fmt"
	"testing"
)

var bms = []struct {
	S    string
	B    []byte
	Want int
}{
	{
		S:    "Firenze;18.9",
		B:    []byte("Firenze;18.9"),
		Want: 7,
	},
	{
		S:    "Las Palmas de Gran Canaria;28.9",
		B:    []byte("Las Palmas de Gran Canaria;28.9"),
		Want: 27,
	},
}

func BenchmarkWithCut(b *testing.B) {
	for _, bm := range bms {
		name := fmt.Sprintf("len-%v", len(bm.S))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				WithCut(bm.S, ";")
			}
		})
	}
}
func BenchmarkWithSplit(b *testing.B) {
	for _, bm := range bms {
		name := fmt.Sprintf("len-%v", len(bm.S))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				WithSplit(bm.S, ";")
			}
		})
	}
}

func BenchmarkWithIndex(b *testing.B) {
	for _, bm := range bms {
		name := fmt.Sprintf("len-%v", len(bm.S))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				WithIndex(bm.S, ";")
			}
		})
	}
}

func BenchmarkWithIter(b *testing.B) {
	for _, bm := range bms {
		name := fmt.Sprintf("len-%v", len(bm.S))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				WithIter(bm.S, ";")
			}
		})
	}
}

func BenchmarkWithSwar(b *testing.B) {
	for _, bm := range bms {
		name := fmt.Sprintf("len-%v", len(bm.S))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				WithSwar(bm.B)
			}
		})
	}
}

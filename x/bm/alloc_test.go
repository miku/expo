package bm

import "testing"

func BenchmarkReserve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reserve()
	}
}

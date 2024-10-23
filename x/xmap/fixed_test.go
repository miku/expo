package xmap

import "testing"

func BenchmarkStaticMap(b *testing.B) {
	m := NewStaticMap()
	key := "hello"
	v := &Measurements{}
	for i := 0; i < b.N; i++ {
		m.Set(key, v)
	}
}

package xmap

import "testing"

func TestXMap(t *testing.T) {
	m := &M{
		NumBuckets: 16,
	}
	m.Init()
	m.Set([]byte("hello"), "go")
	v := m.Get([]byte("hello"))
	t.Logf("got: %v", v)
}

func BenchmarkXMap(b *testing.B) {
	// 131072
	m := &M{
		NumBuckets: 131072,
	}
	m.Init()
	key := []byte("hello")
	for i := 0; i < b.N; i++ {
		m.Set(key, "go")
	}
}

func BenchmarkStdMap(b *testing.B) {
	m := make(map[string]any)
	key := "hello"
	for i := 0; i < b.N; i++ {
		m[key] = "go"
	}
}

package xmap

import "bytes"

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

type Item struct {
	Key   []byte
	Value any
}

type M struct {
	NumBuckets int
	Items      []Item
	size       int
}

func (m *M) Init() {
	m.Items = make([]Item, m.NumBuckets)
}

func (m *M) key(p []byte) uint64 {
	hash := uint64(offset64)
	for _, c := range p {
		hash ^= uint64(c)
		hash *= prime64
	}
	return hash
}

func (m *M) Set(key []byte, value any) {
	u := m.key(key)
	idx := int(u & uint64(m.NumBuckets-1))
	for {
		if m.Items[idx].Key == nil {
			// empty slot, add new item
			k := make([]byte, len(key))
			_ = copy(k, key)
			m.Items[idx] = Item{
				Key:   k,
				Value: value,
			}
			m.size++
			if m.size > m.NumBuckets/2 {
				panic("too many items")
			}
			break
		}
		if bytes.Equal(m.Items[idx].Key, key) {
			// existing entry, update value
			m.Items[idx].Value = value
			break
		}
		idx++
		if idx > m.NumBuckets {
			idx = 0
		}
	}
}

func (m *M) Get(key []byte) any {
	u := m.key(key)
	idx := int(u & uint64(m.NumBuckets-1))
	for {
		if m.Items[idx].Key == nil {
			return nil
		}
		if bytes.Equal(m.Items[idx].Key, key) {
			// existing entry, update value
			return m.Items[idx].Value
		}
		idx++
		if idx > m.NumBuckets {
			idx = 0
		}
	}
	return nil
}

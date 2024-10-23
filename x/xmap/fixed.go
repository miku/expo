package xmap

// Measurements, as there is no need to keep all numbers around, we can compute
// them on the fly.
type Measurements struct {
	Min   int
	Max   int
	Sum   int
	Count int
}

// StaticMap is a tailored "map" for cities.txt.
type StaticMap struct {
	M []*Measurements
}

func NewStaticMap() *StaticMap {
	return &StaticMap{
		M: make([]*Measurements, 16384),
	}
}

// calculateIndex, interestingly the most expensive part of the program.
func calculateIndex(s string) (index int) {
	for i, c := range s {
		index = index + i*(37+int(c))
	}
	return index % 16384
}

func (m *StaticMap) Index(s string) (index int) {
	return calculateIndex(s)
}

func (m *StaticMap) Init() {
	m.M = make([]*Measurements, 16384)
}

func (m *StaticMap) Set(key string, ms *Measurements) {
	m.M[m.Index(key)] = ms
}

func (m *StaticMap) Get(key string) *Measurements {
	return m.M[m.Index(key)]
}

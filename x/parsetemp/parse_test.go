package parsetemp

import "testing"

func BenchmarkParseTempFloat(b *testing.B) {
	var data = "17.3"
	for i := 0; i < b.N; i++ {
		ParseTempFloat(data)
	}
}
func BenchmarkParseTempBigFloat(b *testing.B) {
	var data = "17.3"
	for i := 0; i < b.N; i++ {
		ParseTempBigFloat(data)
	}
}

func BenchmarkParseTempToInt(b *testing.B) {
	var data = []byte("17.3")
	for i := 0; i < b.N; i++ {
		ParseTempToInt(data)
	}
}

func BenchmarkParseTempToIntShift(b *testing.B) {
	var data = []byte("17.3")
	for i := 0; i < b.N; i++ {
		ParseTempToInt(data)
	}
}

func BenchmarkParseNumber(b *testing.B) {
	var data = []byte("17.3")
	for i := 0; i < b.N; i++ {
		ParseNumber(data)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/parsetemp
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// BenchmarkParseTempFloat
// BenchmarkParseTempFloat-8               57792348                20.68 ns/op            0 B/op          0 allocs/op
// BenchmarkParseTempBigFloat
// BenchmarkParseTempBigFloat-8             2588124               482.6 ns/op           200 B/op          8 allocs/op
// BenchmarkParseTempToInt
// BenchmarkParseTempToInt-8               301726612                4.013 ns/op           0 B/op          0 allocs/op
// BenchmarkParseTempToIntShift
// BenchmarkParseTempToIntShift-8          295691558                3.999 ns/op           0 B/op          0 allocs/op
// BenchmarkParseNumber
// BenchmarkParseNumber-8                  1000000000               0.2538 ns/op          0 B/op          0 allocs/op

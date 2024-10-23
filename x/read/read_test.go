package read

import (
	"fmt"
	"testing"
)

func BenchmarkRaw(b *testing.B) {
	var bms = []struct {
		Filename string
		Size     int
	}{
		{"../../measurements.txt", 1024},
		{"../../measurements.txt", 2048},
		{"../../measurements.txt", 4096},
		{"../../measurements.txt", 8192},
		{"../../measurements.txt", 16384},
		{"../../measurements.txt", 32768},
		{"../../measurements.txt", 65536},
		{"../../measurements.txt", 131072},
		{"../../measurements.txt", 262144},
		{"../../measurements.txt", 524288},
		{"../../measurements.txt", 1048576},
		{"../../measurements.txt", 2097152},
	}
	for _, bm := range bms {
		name := fmt.Sprintf("size-%d", bm.Size)
		b.Run(name, func(b *testing.B) {
			if err := Raw(bm.Filename, bm.Size); err != nil {
				b.Fatalf("raw: %v", err)
			}
		})
	}
	// goos: linux
	// goarch: amd64
	// pkg: github.com/miku/expo/x/read
	// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
	// PASS
	// benchmark                     iter        time/iter
	// ---------                     ----        ---------
	// BenchmarkRaw/size-1024-8         1   10826.98 ms/op
	// BenchmarkRaw/size-2048-8         1    8715.33 ms/op
	// BenchmarkRaw/size-4096-8         1    7463.78 ms/op
	// BenchmarkRaw/size-8192-8         1    6767.17 ms/op
	// BenchmarkRaw/size-16384-8        1    6743.58 ms/op
	// BenchmarkRaw/size-32768-8        1    6588.15 ms/op
	// BenchmarkRaw/size-65536-8        1    6372.79 ms/op
	// BenchmarkRaw/size-131072-8       1    6302.75 ms/op
	// BenchmarkRaw/size-262144-8       1    6057.32 ms/op
	// BenchmarkRaw/size-524288-8       1    5921.70 ms/op
	// BenchmarkRaw/size-1048576-8      1    6129.24 ms/op
	// BenchmarkRaw/size-2097152-8      1    7229.02 ms/op
	// ok      github.com/miku/expo/x/read     85.130s
}

func BenchmarkReadByte(b *testing.B) {
	var bms = []struct {
		Filename string
	}{
		{"../../measurements.txt"},
	}
	for _, bm := range bms {
		b.Run("m", func(b *testing.B) {
			if err := ReadByte(bm.Filename); err != nil {
				b.Fatalf("read byte: %v", err)
			}
		})
	}
}

func BenchmarkReadString(b *testing.B) {
	var bms = []struct {
		Filename string
	}{
		{"../../measurements.txt"},
	}
	for _, bm := range bms {
		b.Run("m", func(b *testing.B) {
			if err := ReadString(bm.Filename); err != nil {
				b.Fatalf("read string: %v", err)
			}
		})
	}
}

func BenchmarkReadBytes(b *testing.B) {
	var bms = []struct {
		Filename string
	}{
		{"../../measurements.txt"},
	}
	for _, bm := range bms {
		b.Run("m", func(b *testing.B) {
			if err := ReadBytes(bm.Filename); err != nil {
				b.Fatalf("read bytes: %v", err)
			}
		})
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/read
// cpu: 13th Gen Intel(R) Core(TM) i9-13900T
// PASS
// benchmark                  iter    time/iter
// ---------                  ----    ---------
// BenchmarkReadByte/m-32        1   23.22 s/op
// BenchmarkReadString/m-32      1   37.64 s/op
// BenchmarkReadBytes/m-32       1   36.01 s/op

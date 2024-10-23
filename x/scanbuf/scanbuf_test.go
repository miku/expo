package scanbuf

import (
	"fmt"
	"testing"
)

func BenchmarkWithBuffer(b *testing.B) {
	var bms = []struct {
		Filename string
		Size     int
	}{
		{"../../measurements.txt", 1024},
		{"../../measurements.txt", 4096},
		{"../../measurements.txt", 16384},
		{"../../measurements.txt", 65536},
		{"../../measurements.txt", 262144},
		{"../../measurements.txt", 524288},
		{"../../measurements.txt", 1048576},
		{"../../measurements.txt", 2097152},
		{"../../measurements.txt", 4194304},
		{"../../measurements.txt", 8388608},
		{"../../measurements.txt", 16777216},
		{"../../measurements.txt", 33554432},
	}
	for _, bm := range bms {
		name := fmt.Sprintf("buf-%d", bm.Size)
		b.Run(name, func(b *testing.B) {
			err := WithBuffer(bm.Filename, bm.Size)
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
		})
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/scanbuf
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
//
// PASS
// benchmark                            iter    time/iter
// ---------                            ----    ---------
// BenchmarkWithBuffer/buf-1024-8          1   34.03 s/op
// BenchmarkWithBuffer/buf-4096-8          1   23.15 s/op
// BenchmarkWithBuffer/buf-16384-8         1   21.81 s/op
// BenchmarkWithBuffer/buf-65536-8         1   20.96 s/op
// BenchmarkWithBuffer/buf-262144-8        1   23.26 s/op
// BenchmarkWithBuffer/buf-524288-8        1   20.38 s/op
// BenchmarkWithBuffer/buf-1048576-8       1   22.37 s/op
// BenchmarkWithBuffer/buf-2097152-8       1   23.89 s/op
// BenchmarkWithBuffer/buf-4194304-8       1   24.87 s/op
// BenchmarkWithBuffer/buf-8388608-8       1   24.81 s/op
// BenchmarkWithBuffer/buf-16777216-8      1   25.24 s/op
// BenchmarkWithBuffer/buf-33554432-8      1   25.66 s/op
// ok      github.com/miku/expo/x/scanbuf  290.446s
//
// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/scanbuf
// cpu: 13th Gen Intel(R) Core(TM) i9-13900T
// PASS
// benchmark                             iter        time/iter
// ---------                             ----        ---------
// BenchmarkWithBuffer/buf-1024-32          1   14209.42 ms/op
// BenchmarkWithBuffer/buf-4096-32          1   10123.61 ms/op
// BenchmarkWithBuffer/buf-16384-32         1    9773.66 ms/op
// BenchmarkWithBuffer/buf-65536-32         1    8797.50 ms/op
// BenchmarkWithBuffer/buf-262144-32        1    8413.58 ms/op
// BenchmarkWithBuffer/buf-524288-32        1    9169.79 ms/op
// BenchmarkWithBuffer/buf-1048576-32       1    8542.06 ms/op
// BenchmarkWithBuffer/buf-2097152-32       1    8626.86 ms/op
// BenchmarkWithBuffer/buf-4194304-32       1    9280.23 ms/op
// BenchmarkWithBuffer/buf-8388608-32       1    8476.03 ms/op
// BenchmarkWithBuffer/buf-16777216-32      1    8763.91 ms/op
// BenchmarkWithBuffer/buf-33554432-32      1    9709.72 ms/op

package split

import (
	"bytes"
	"strings"
)

var bSep = []byte{';'}

// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/split
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkStringSplit/s-len-9-8          26202700                65.26 ns/op
// BenchmarkStringSplit/s-len-31-8         16804399                61.59 ns/op
// BenchmarkByteSliceSplit/b-len-9-8       16961588                91.13 ns/op
// BenchmarkByteSliceSplit/b-len-31-8      19413685                59.67 ns/op
// BenchmarkCustomSplit/b-len-9-8          332549199                3.603 ns/op
// BenchmarkCustomSplit/b-len-31-8         151124314                7.992 ns/op

func StringSplit(s string) any {
	fields := strings.Split(s, ";")
	return fields
}

func BytesSplit(b []byte) any {
	var bSep = []byte{';'}
	fields := bytes.Split(b, bSep)
	return fields
}

func CustomSplit(p []byte) any {
	var fields = make([][]byte, 2)
	for i, c := range p {
		if c == ';' {
			fields[0] = p[0:i]
			fields[1] = p[i : len(p)-2]
			break
		}
	}
	return fields
}

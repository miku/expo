# Basic

Benchmark time and allocations.

```
$ make
go test -v -bench=. -benchmem
=== RUN   TestFindContains
--- PASS: TestFindContains (0.00s)
=== RUN   TestFindIter
--- PASS: TestFindIter (0.00s)
goos: linux
goarch: amd64
pkg: github.com/miku/expo/x/bm
cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
BenchmarkReserve
BenchmarkReserve-8                114326             11135 ns/op           82664 B/op         15 allocs/op
BenchmarkFindContains
BenchmarkFindContains/len-27
BenchmarkFindContains/len-27-8          149409480                7.601 ns/op           0 B/op          0 allocs/op
BenchmarkFindContains/len-21
BenchmarkFindContains/len-21-8          98419642                12.89 ns/op            0 B/op          0 allocs/op
BenchmarkFindContains/len-163864
BenchmarkFindContains/len-163864-8        314902              3784 ns/op               0 B/op          0 allocs/op
BenchmarkFindIter
BenchmarkFindIter/len-27
BenchmarkFindIter/len-27-8              378540080                3.201 ns/op           0 B/op          0 allocs/op
BenchmarkFindIter/len-21
BenchmarkFindIter/len-21-8              62152454                19.42 ns/op            0 B/op          0 allocs/op
BenchmarkFindIter/len-163864
BenchmarkFindIter/len-163864-8              9368            126495 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/miku/expo/x/bm       10.781s
```

An older (now archived) project called [prettybench](https://github.com/cespare/prettybench):

```
$ make | prettybench
go test -v -bench=. -benchmem
=== RUN   TestFindContains
--- PASS: TestFindContains (0.00s)
=== RUN   TestFindIter
--- PASS: TestFindIter (0.00s)
goos: linux
goarch: amd64
pkg: github.com/miku/expo/x/bm
cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
BenchmarkReserve
BenchmarkFindContains
BenchmarkFindContains/len-27
BenchmarkFindContains/len-21
BenchmarkFindContains/len-163864
BenchmarkFindIter
BenchmarkFindIter/len-27
BenchmarkFindIter/len-21
BenchmarkFindIter/len-163864
PASS
benchmark                                 iter         time/iter   bytes alloc         allocs
---------                                 ----         ---------   -----------         ------
BenchmarkReserve-8                      113184    11553.00 ns/op    82664 B/op   15 allocs/op
BenchmarkFindContains/len-27-8       149598163        7.70 ns/op        0 B/op    0 allocs/op
BenchmarkFindContains/len-21-8        97313648       13.13 ns/op        0 B/op    0 allocs/op
BenchmarkFindContains/len-163864-8      305828     3816.00 ns/op        0 B/op    0 allocs/op
BenchmarkFindIter/len-27-8           372878384        3.23 ns/op        0 B/op    0 allocs/op
BenchmarkFindIter/len-21-8            58729078       19.83 ns/op        0 B/op    0 allocs/op
BenchmarkFindIter/len-163864-8            9327   127716.00 ns/op        0 B/op    0 allocs/op
ok      github.com/miku/expo/x/bm       9.802s

```

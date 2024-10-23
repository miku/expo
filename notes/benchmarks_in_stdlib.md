# benchmarks in the standard library

Over 1700 benchmarks in the standard library (how many SLOC, and what
percentage)? Across 255 files.

There are 10184 Go files in the distribution (Go 1.23.2), and 255 contain
benchmarks; so about 2.5%.

```
tir@reka:~/code/golang/go [git:master] $ rg -H 'func Benchmark.*(b \*testing.B)' | grep ^src | wc -l
1770
```

Benchmarks live in test code.

```
tir@reka:~/code/golang/go [git:master] $ rg -H 'func Benchmark.*(b \*testing.B)' | grep ^src | cut -d : -f 1 | sort | uniq -c | sort -nr | head -40
     72 src/math/all_test.go
     68 src/runtime/memmove_test.go
     46 src/strings/strings_test.go
     46 src/runtime/map_benchmark_test.go
     43 src/math/bits/bits_test.go
     43 src/bytes/bytes_test.go
     36 src/cmd/compile/internal/test/divconst_test.go
     31 src/unicode/utf8/utf8_test.go
     31 src/encoding/binary/binary_test.go
     29 src/math/rand/v2/rand_test.go
     28 src/time/time_test.go
     27 src/runtime/iface_test.go
     27 src/reflect/benchmark_test.go
     25 src/regexp/all_test.go
     25 src/encoding/json/bench_test.go
     24 src/sync/map_bench_test.go
     22 src/runtime/chan_test.go
     22 src/math/cmplx/cmath_test.go
     22 src/fmt/fmt_test.go
     20 src/runtime/vlop_arm_test.go
     20 src/image/draw/bench_test.go
     20 src/cmd/compile/internal/test/math_test.go
     19 src/net/netip/netip_test.go
     19 src/internal/runtime/atomic/bench_test.go
     18 src/sort/sort_test.go
     18 src/image/image_test.go
     17 src/math/big/int_test.go
     16 src/strings/replace_test.go
     16 src/net/http/serve_test.go
     16 src/math/rand/rand_test.go
     16 src/expvar/expvar_test.go
     16 src/encoding/gob/timing_test.go
     15 src/math/big/gcd_test.go
     13 src/runtime/slice_test.go
     13 src/runtime/runtime_test.go
     13 src/runtime/pinner_test.go
     12 src/strconv/atof_test.go
     12 src/crypto/md5/md5_test.go
     11 src/runtime/string_test.go
     11 src/runtime/proc_test.go
```

# Allocations

Count allocations during tests.

```shell
$ go test -bench=. -benchmem
                   ---------
```

What does `benchmem` measure?

[runtime.MemStats](https://pkg.go.dev/runtime#MemStats) gives insight into the
runtime memory management.

> There are several memory metrics. In Go, some of the most useful are HeapSys
> and HeapAlloc. The first indicates how much memory (in bytes) has been given
> to the program by the operating system. The second value, which is typically
> lower indicates how much of that memory is actively in used by the program.

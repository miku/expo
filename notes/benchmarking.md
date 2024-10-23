# Benchmarking

Similar to writing a test, part of [package testing](https://pkg.go.dev/testing#hdr-Benchmarks).

```go
func BenchmarkRandInt(b *testing.B) {
    for range b.N {
        rand.Int()
    }
}
```

What is [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)?

> Benchstat computes statistical summaries and A/B comparisons of Go benchmarks.

## Memory in Go

> In Go, each function has its own ‘stack memory’. As the name suggests, stack
memory is allocated and deallocated in a last-in, first-out (LIFO) order. This
memory is typically only usable within the function, and it is often limited in
size. The other type of memory that a Go program may use is heap memory. Heap
memory is allocated and deallocated in a random order. There is only one heap
shared by all functions. -- [Measuring memory allocations](https://lemire.me/blog/2024/03/17/measuring-your-systems-performance-using-software-go-edition/)

And following:

> Go automatically reclaims it. However, it may still be bad for performance to
> constantly allocate and deallocate memory. In many real-world systems, memory
> management becomes a performance bottleneck.

Rule of thumb:

> Typically, in Go, it roughly corresponds to the number of calls to make and
> to the number of objects that the garbage collector must handle.

## Escape analysis

* Go refspec does not mention heap or stack at all
* a value may live on the stack or heap, depending on a runtime
* a value that looks heap allocated may live in the stack (because it fits and is not used outside of the function)
* a value that looks stack allocated may live on the heap, because it is used after the function returned

## Page size

What is the page size? Smallest unit of memory.

> Your operating system provides memory to a running process in units of pages.
> The operating system cannot provide memory in smaller units than a page.

```
$ getconf PAGESIZE
```

It is difficult to get exact usage amounts.

> Thus it is not simple to ask how much memory a program uses. A program may
> appear to use a lot of (virtual) memory, while not using much physical
> memory, and inversely.

## Cost of allocation

Kernel overhead, must zero the memory. Bookkeeping.

> Allocating pages to a process is not free, it takes some effort. Among other
> things, the operating system cannot just reuse a memory page from another
> process as is. Doing so would be a security threat because you could have
> indirect access to the data stored in memory by another process. This other
> process could have held in memory your passwords or other sensitive
> information. Typically an operating system has to initialize (e.g., set to
> zero) a newly assigned page.

## Tips

* [https://llvm.org/docs/Benchmarking.html](https://llvm.org/docs/Benchmarking.html)



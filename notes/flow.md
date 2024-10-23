# flow

* Intro
    * Inspiration
        * Exercises in Style (1960)
        * Exercises in Programming Style
          ([book](https://www.routledge.com/Exercises-in-Programming-Style/Lopes/p/book/9780367350208),
           [repo](https://github.com/crista/exercises-in-programming-style),
           [talk](https://www.infoq.com/presentations/programming-styles/))
    * While not constrained programming or writing, we look at variations of a single problem

* Benchmarking
    * a simple benchmark
        * testing.B
        * using sub-benchmarks

    ```go
    // BenchmarkResult contains the results of a benchmark run.
    type BenchmarkResult struct {
        N         int           // The number of iterations.
        T         time.Duration // The total time taken.
        Bytes     int64         // Bytes processed in one iteration.
        MemAllocs uint64        // The total number of memory allocations.
        MemBytes  uint64        // The total number of bytes allocated.

        // Extra records additional metrics reported by ReportMetric.
        Extra map[string]float64
    }
    ```

    * running a benchmark
        * allocations from `runtime.Memstats`
    * output is defined, in [14313](https://go.googlesource.com/proposal/+/master/design/14313-benchmark-format.md)

> Sidenote: The Go team discovers these small hindrances and tries to fix them.
> Cf. "We are unaware of any standard formats for recording raw benchmark data,
> and we've been unable to find any using web searches. One might expect that a
> standard benchmark suite such as SPEC CPU2006 would have defined a format for
> raw results, but that appears not to be the case."


* Profiling
    * To find hot spots, enable cpu and memory profiling
        * Example snippet
    * What to do with a profile?
        * export to PNG
        * explore with web
        * various pprof profiles, goroutine, heap, threadcreate, block, mutex, cpu and trace
    * The pprof format seems to be a Google side project;
      https://github.com/google/pprof based on
      [profile.proto](https://github.com/google/pprof/blob/main/proto/profile.proto)

    > Sidenote, you can use
    > [perf_data_converter](https://github.com/google/perf_data_converter) and use
    > pprof with output from Linux perf tool. TODO: show an example

* Additional tools
    * https://pkg.go.dev/golang.org/x/perf#section-readme
    * https://pkg.go.dev/github.com/aclements/go-misc/benchmany

> benchmany runs the benchmarks in the current directory <iterations> times for
> each commit in <commit or range> and writes the benchmark results to
> bench.log. Benchmarks may be Go testing framework benchmarks or benchmarks
> from golang.org/x/benchmarks.

This in turn allows you to find performance degradations.

    * also: https://github.com/aclements/go-misc/tree/master/benchplot
    * TODO: what is? https://perf.golang.org/

Convert pprof to Brendan Gregg format: https://github.com/felixge/pprofutils

> pprofutils is a swiss army knife for pprof files. You can use it as a command line utility or as a free web service.

* Flame Graphs
    * https://queue.acm.org/detail.cfm?id=2927301, The Flame Graph

> An everyday problem in our industry is understanding how software is
> consuming resources, particularly CPUs. What exactly is consuming how much,
> and how did this change since the last software version? These questions can
> be answered using software profilers, tools that help direct developers to
> optimize their code and operators to tune their environment. The output of
> profilers can be verbose, however, making it laborious to study and
> comprehend. The flame graph provides a new visualization for profiler output
> and can make for much faster comprehension, reducing the time for root cause
> analysis.




## Raw

> About 120 min, 12x10

1. run benchmarking w/ testing.B
    * run one benchmark
    * what happens under the hood, memstats
2. tools: pprof (png, web), perf
    * pprof format
    * pretty output
3. baselines
    * plain iteration, wc, cw, cat
    * bufio read vs scanner
    * plain scanner
    * complete data in memory
4. scanner optimizations
    * tweaking a scanner
5. parallelize
    * parallelize with basic pattern
    * batching, transfer of ownership
        * benchmark various options and batch sizes, separately
6. mmap
    * using mmap to access file, https://unix.stackexchange.com/q/474926/376
    * mmap on various platforms
7. float ops
    * parsing a float
    * writing a float (strconv, formatfloat, fmt functions)
    * just use int
    * benchmark various way to parse/read floats
        * https://stackoverflow.com/a/55326309/89391
8. use a custom map
    * how fast is the builtin map
    * are there other implementations?
    * building a custom (perfect) hash map
9. swar
    * examples

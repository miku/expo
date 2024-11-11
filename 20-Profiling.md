# Profiling

Go has a few resource profiler:

* cpu
* memory
* goroutines
* heap
* thread (os)
* block
* mutex

Can be used with http services via `import _ "net/http/pprof"`, or for plain programs.

```go
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
...
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
...
```

The output is in the
[profile.proto](https://github.com/google/pprof/blob/main/proto/profile.proto)
format, that [pprof](https://github.com/google/pprof) can read.

## go tool pprof

The pprof tool subcommand allows to render profiling files in the browser or to an SVG.

### Rendering a pprof file

To the web:

```
$ go tool pprof -http=: cpu.pprof
```

To an svg:

```
$ go tool pprof -png cpu.pprof > cpu.png

```

Cf. [git.io/JfYMW](git.io/JfYMW)

## Flamegraphs

Via web interface, or via standalone tools (scripts).

E.g. via [brendangregg/FlameGraph](https://github.com/brendangregg/FlameGraph)

```shell
#!/bin/bash

PPROF=${1:-cpu.pprof}
SVG=${2:-cpu.svg}

go tool pprof -raw -output=cpu.txt "$PPROF"
stackcollapse-go.pl cpu.txt | flamegraph.pl > "$SVG"
```

## perf (linux)

> Perf is a profiler tool for Linux 2.6+ based systems that abstracts away CPU
> hardware differences in Linux performance measurements and presents a simple
> commandline interface. Perf is based on the perf_events interface exported by
> recent versions of the Linux kernel.

```
$ perf stat -e branches,branch-misses,cache-references,cache-misses,cycles,instructions -- you_program ...
```

Example:

```
$ perf stat -e branches,branch-misses,cache-references,cache-misses,cycles,instructions -- wc -l measurements.txt
1000000000 measurements.txt

 Performance counter stats for 'wc -l measurements.txt':

     3,536,931,401      branches                                                                (83.01%)
        15,962,856      branch-misses                    #    0.45% of all branches             (83.61%)
       908,243,138      cache-references                                                        (83.63%)
       441,657,605      cache-misses                     #   48.63% of all cache refs           (83.05%)
    20,047,628,276      cycles                                                                  (66.70%)
    21,228,533,559      instructions                     #    1.06  insn per cycle              (83.39%)

       9.259674464 seconds time elapsed

       0.539239000 seconds user
       5.314362000 seconds sys
```

Note: On MacOS there may be an issue with generating CPU profiles, like it
[overcounts system calls](https://github.com/golang/go/issues/57722).


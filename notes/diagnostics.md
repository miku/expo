# Diagnostics

* [https://go.dev/doc/diagnostics](https://go.dev/doc/diagnostics)

> The Go ecosystem provides a large suite of APIs and tools to diagnose logic
> and performance problems in Go programs. This page summarizes the available
> tools and helps Go users pick the right one for their specific problem.

* Profiling: Profiling tools analyze the complexity and costs of a Go program
  such as its memory usage and frequently called functions to identify the
  expensive sections of a Go program.
* Tracing: Tracing is a way to instrument code to analyze latency throughout
  the lifecycle of a call or user request. Traces provide an overview of how
  much latency each component contributes to the overall latency in a system.
  Traces can span multiple Go processes.
* Debugging: Debugging allows us to pause a Go program and examine its
  execution. Program state and flow can be verified with debugging.
* Runtime statistics and events: Collection and analysis of runtime stats and
  events provides a high-level overview of the health of Go programs.
  Spikes/dips of metrics helps us to identify changes in throughput, utilization,
  and performance.

What can we use:

* benchmarking with `testing.B`
* shell: `time.time`
* profiling of Go code, to identify potential issues

Using Linux tools, like [perf](https://perfwiki.github.io/main/).

> Performance counters are CPU hardware registers that count hardware events
> such as instructions executed, cache-misses suffered, or branches
> mispredicted. They form a basis for profiling applications to trace dynamic
> control flow and identify hotspots. perf provides rich generalized
> abstractions over hardware specific capabilities. Among others, it provides
> per task, per CPU and per-workload counters, sampling on top of these and
> source code event annotation.

# perf (linux)

```
$ cat /proc/sys/kernel/perf_event_paranoid
4
```

* -1: Allow use of (almost) all events by all users
      Ignore mlock limit after `perf_event_mlock_kb` without `CAP_IPC_LOCK`
* >=0: Disallow ftrace function tracepoint by users without `CAP_SYS_ADMIN`
       Disallow raw tracepoint access by users without `CAP_SYS_ADMIN`
* >=1: Disallow CPU event access by users without `CAP_SYS_ADMIN`
* >=2: Disallow kernel profiling by users without `CAP_SYS_ADMIN`

Two modes.

> As a profiler, perf provides two working mode: sampling (`perf record`) and
  counting (`perf stat`).

Allow tracing.

```
echo 2 | sudo tee /proc/sys/kernel/perf_event_paranoid
```

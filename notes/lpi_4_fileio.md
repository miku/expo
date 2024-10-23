# file io

* all io syscalls use fds
* pipes, FIFOs, sockets, terminals, devices, and regular files
* each process has its own set of file descriptors

```
0: stdin, STDIN_FILENO
1: stdout, STDOUT_FILENO
2: stderr, STDERR_FILENO
```

From [unistd.h](https://en.wikipedia.org/wiki/Unistd.h)

FDs can be reopened with [freopen](https://pubs.opengroup.org/onlinepubs/007904875/functions/freopen.html)

* open
* read, is using a buffer as well
* write
* close

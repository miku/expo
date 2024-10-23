# Memory-mapped file

> mmap() creates a new mapping in the virtual address space of the calling process. -- [mmap(2)](https://man7.org/linux/man-pages/man2/mmap.2.html)

Advantages:

* useful for parallel access to read-only data; inter-process communication (IPC)
* no lseek, generally faster

Disadvantages:

* size must be multiple of pagesize (e.g. 4KB)
* may be less useful with smaller files, as there is initial overhead in setting up memory mapping

Support in Go with x/exp:

* [https://pkg.go.dev/golang.org/x/exp/mmap](https://pkg.go.dev/golang.org/x/exp/mmap)

Implements [io.ReaderAt](https://pkg.go.dev/io#ReaderAt).

> Like any io.ReaderAt, clients can execute parallel ReadAt calls

Lazy I/O.

## Misc

* [So what’s wrong with 1975 programming ?](https://varnish-cache.org/docs/trunk/phk/notes.html)

> The really short answer is that computers do not have two kinds of storage anymore.

> It used to be that you had the primary store, and it was anything from acoustic delaylines filled with mercury via small magnetic doughnuts via transistor flip-flops to dynamic RAM.

> And then there were the secondary store, paper tape, magnetic tape, disk drives the size of houses, then the size of washing machines and these days so small that girls get disappointed if think they got hold of something else than the MP3 player you had in your pocket.

> And people program this way.

> They have variables in “memory” and move data to and from “disk”.
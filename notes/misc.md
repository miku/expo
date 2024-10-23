# misc

* [How to improve the speed of reading a large file line by line in Go](https://stackoverflow.com/questions/55010716/how-to-improve-the-speed-of-reading-a-large-file-line-by-line-in-go)
* [What does I/O bound really mean?](http://erikengbrecht.blogspot.com/2008/06/what-does-io-bound-really-mean.html)

Reading a large file.

* [Leveraging multithreading to read large files faster in Go](https://medium.com/@mohdgadi52/leveraging-multithreading-to-read-large-files-faster-in-go-cfb9d6a77aeb), missing `io.ReaderAt`

Line by line options:

* [Reading a File Line by Line in Go: A Comprehensive Guide](https://www.bacancytechnology.com/qanda/golang/reading-a-file-line-by-line-in-go)

Reading files:

* [Reading files in Go — an overview](https://kgrz.io/reading-files-in-go-an-overview.html)

The [read](https://www.man7.org/linux/man-pages/man2/read.2.html#NOTES) syscall:

> On Linux, read() (and similar system calls) will transfer at most 0x7ffff000
> (2,147,479,552) bytes, returning the number of bytes actually transferred.
> (This is true on both 32-bit and 64-bit systems.)

## unexplored optimization techniques

* using `io_uring` w/ Go
* simd
* read, pread, preadv

The calls: https://man7.org/linux/man-pages/man2/readv.2.html

## Perfect hash map

* https://homepages.dcc.ufmg.br/~nivio/papers/cikm07.pdf
* http://burtleburtle.net/bob/hash/perfect.html

CHD, but its a read only thing; https://github.com/alecthomas/mph?tab=readme-ov-file#what-is-this-useful-for

## The numbers

* range: -99.9 to 99.9
* int variant: -999 to 999

## Maps

* avoid map growth by preallocating space for e.g. 413 cities

## String interning

We have 413 unique keys, we do not need new objects all the time: https://artem.krylysov.com/blog/2018/12/12/string-interning-in-go/

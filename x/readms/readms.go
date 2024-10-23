package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

func main() {
	ms := new(runtime.MemStats)
	fmt.Printf("size  ms: %v\n", unsafe.Sizeof(ms))
	fmt.Printf("size *ms: %v\n", unsafe.Sizeof(*ms))
	runtime.ReadMemStats(ms)
	fmt.Printf("heap alloc=%v num allocs=%v\n", ms.HeapAlloc, ms.Mallocs)
	runtime.ReadMemStats(ms)
	fmt.Printf("heap alloc=%v num allocs=%v\n", ms.HeapAlloc, ms.Mallocs)
}

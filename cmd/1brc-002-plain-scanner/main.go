// How fast is a plain scanner over a file? See x/scanbuf/ for a benchmark
// comparing various buffer sizes.
//
// i7-8550U
//
// scanner.Text()
//
// real    0m27.617s
// user    0m19.938s
// sys     0m7.200s
//
// scanner.Bytes()
//
// real    0m22.092s
// user    0m14.010s
// sys     0m7.104s

package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		_ = scanner.Text() // allocate
		//_ = scanner.Bytes() // does not allocate
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

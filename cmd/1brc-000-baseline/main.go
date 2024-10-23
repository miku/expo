// Keep all data in memory and perform stats on slice. Requires O(n) memory.
//
// data:
//
// Tamale;27.5
// Bergen;9.6
// Lodwar;37.1
// Whitehorse;-3.8
// Ouarzazate;19.1
//
// i7-8550U
//
// real    3m52.862s
// user    3m42.343s
// sys     0m11.108s
//
// i9-13900T
//
// real    2m31.412s
// user    2m27.340s
// sys     0m6.374s

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

var data = make(map[string][]float32)

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
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			log.Fatalf("expected two fields: %s", line)
		}
		name := strings.TrimSpace(parts[0])
		temp, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			log.Fatalf("invalid temp: %s", parts[1])
		}
		data[name] = append(data[name], float32(temp))
	}
	keys := maps.Keys(data)
	sort.Strings(keys)
	for _, k := range keys {
		min := slices.Min(data[k])
		max := slices.Max(data[k])
		var sum float32 = 0.0
		for _, t := range data[k] {
			sum = sum + t
		}
		avg := sum / float32(len(data[k]))
		fmt.Printf("%s\t%0.2f/%0.2f/%0.2f\n", k, min, max, avg)
	}
}

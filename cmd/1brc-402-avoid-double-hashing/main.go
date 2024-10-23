// R2: Avoid double hashing.
//
// R1: stats[name] was access twice, hence calculating a hash, twice.
//
// # Intel i7-8550U
//
// real    2m7.445s
// user    2m1.394s
// sys     0m6.043s
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"runtime/pprof"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

type Stats struct {
	Min, Max, Sum float64
	Count         int64
}

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
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	if err := R2(*filename, bw); err != nil {
		log.Fatal(err)
	}
}

func R2(fn string, w io.Writer) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	stats := make(map[string]*Stats) // vs R1: value is now a pointer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		name, tempStr, found := strings.Cut(line, ";")
		if !found {
			continue
		}
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return err
		}
		s, ok := stats[name]
		if !ok {
			stats[name] = &Stats{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		} else {
			s.Min = min(s.Min, temp)
			s.Max = max(s.Min, temp)
			s.Sum = s.Sum + temp
			s.Count++
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	names := slices.Collect(maps.Keys(stats))
	sort.Strings(names)
	_, _ = fmt.Fprintf(w, "{")
	for i, name := range names {
		if i > 0 {
			_, _ = fmt.Fprintf(w, ", ")
		}
		s := stats[name]
		mean := s.Sum / float64(s.Count)
		_, _ = fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f", name, s.Min, mean, s.Max)
	}
	fmt.Fprintln(w, "}\n")
	return nil
}

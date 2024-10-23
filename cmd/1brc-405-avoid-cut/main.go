// R5: avoid cut
//
// # Intel i7-8550U
//
// real    1m32.339s
// user    1m26.543s
// sys     0m6.498s
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
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

type Stats struct {
	Min, Max, Sum int32
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
	if err := R5(*filename, bw); err != nil {
		log.Fatal(err)
	}
}

func R5(fn string, w io.Writer) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	stats := make(map[string]*Stats) // vs R1: value is now a pointer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Improve on "Cut"
		// ----8<----
		// name, tempStr, found := strings.Cut(line, ";")
		// if !found {
		// 	continue
		// }
		//
		// Note: this requires the exact format, with a single decimal digit.
		end := len(line)
		tenths := int32(line[end-1] - '0')
		ones := int32(line[end-3] - '0') // line[end-2] is '.'
		var temp int32
		var semicolon int
		if line[end-4] == ';' {
			temp = ones*10 + tenths
			semicolon = end - 4
		} else if line[end-4] == '-' {
			temp = -(ones*10 + tenths)
			semicolon = end - 5
		} else {
			tens := int32(line[end-4] - '0')
			if line[end-5] == ';' {
				temp = tens*100 + ones*10 + tenths
				semicolon = end - 5
			} else {
				temp = -(tens*100 + ones*10 + tenths)
				semicolon = end - 6
			}
		}
		name := line[:semicolon]
		// ----8<----

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
		mean := float64(s.Sum) / float64(s.Count)
		_, _ = fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f",
			name, float64(s.Min)/10, mean, float64(s.Max)/10)
	}
	fmt.Fprintln(w, "}\n")
	return nil
}

// R6: avoid scanner
//
// # Intel i7-8550U
//
// real    0m57.534s
// user    0m51.778s
// sys     0m3.913s
package main

import (
	"bufio"
	"bytes"
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
	if err := R6(*filename, bw); err != nil {
		log.Fatal(err)
	}
}

func R6(fn string, w io.Writer) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	stats := make(map[string]*Stats) // vs R1: value is now a pointer
	// ----8<----
	// do not use bufio.Scanner
	buf := make([]byte, 1<<20)
	readStart := 0
	// ----8<----
	for {
		n, err := f.Read(buf[readStart:])
		if err != nil && err != io.EOF {
			return err
		}
		if readStart+n == 0 {
			break
		}
		chunk := buf[:readStart+n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			break
		}
		remaining := chunk[newline+1:]
		chunk = chunk[:newline+1]

		for {
			name, after, found := bytes.Cut(chunk, []byte(";"))
			if !found {
				break
			}
			i := 0
			neg := false
			if after[i] == '-' {
				neg = true
				i++
			}
			temp := int32(after[i] - '0')
			i++
			if after[i] != '.' {
				temp = temp*10 + int32(after[i]-'0')
				i++
			}
			i++ // skip '.'
			temp = temp*10 + int32(after[i]-'0')
			i += 2 // skip last digit and '\n'
			if neg {
				temp = -temp
			}
			chunk = after[i:]

			s, ok := stats[string(name)]
			if !ok {
				stats[string(name)] = &Stats{
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
		readStart = copy(buf, remaining)
	}
	names := slices.Collect(maps.Keys(stats))
	sort.Strings(names)
	_, _ = fmt.Fprintf(w, "{")
	for i, name := range names {
		if i > 0 {
			_, _ = fmt.Fprintf(w, ", ")
		}
		s := stats[name]
		mean := float64(s.Sum) / float64(s.Count) / 10
		_, _ = fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f",
			name, float64(s.Min)/10, mean, float64(s.Max)/10)
	}
	fmt.Fprintln(w, "}\n")
	return nil
}

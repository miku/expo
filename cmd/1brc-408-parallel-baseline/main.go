// R8: Parallel baseline
//
// CPU: Intel i7-8550U
//
// real    0m55.701s
// user    6m42.195s
// sys     0m13.163s
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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

type Part struct {
	Offset int64
	Size   int64
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
	if err := R8(*filename, bw); err != nil {
		log.Fatal(err)
	}
}

func R8(fn string, w io.Writer) error {
	parts, err := splitFile(fn, runtime.NumCPU())
	if err != nil {
		return err
	}
	// run one goroutine per part
	resultC := make(chan map[string]Stats)
	for _, part := range parts {
		go processPart(fn, part, resultC) // XXX: cancellation and wait?
	}

	// merge stats from goroutines in totals
	totals := make(map[string]Stats)
	for i := 0; i < len(parts); i++ {
		result := <-resultC
		for name, s := range result {
			ts, ok := totals[name]
			if !ok {
				totals[name] = Stats{
					Min:   s.Min,
					Max:   s.Max,
					Sum:   s.Sum,
					Count: s.Count,
				}
			} else {
				ts.Min = min(ts.Min, s.Min)
				ts.Max = max(ts.Max, s.Max)
				ts.Sum += s.Sum
				ts.Count += s.Count
				totals[name] = ts
			}
		}
	}

	// Get the names out.
	names := make([]string, 0, len(totals))
	for station := range totals {
		names = append(names, station)
	}
	sort.Strings(names)

	fmt.Fprint(w, "{")
	for i, name := range names {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		s := totals[name]
		mean := s.Sum / float64(s.Count)
		fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f", name, s.Min, mean, s.Max)
	}
	fmt.Fprint(w, "}\n")
	return nil
}

// processPart does not do any error handling for the moment.
func processPart(fn string, part Part, resultC chan map[string]Stats) {
	// XXX: can we just reuse the file handle?
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Seek(part.Offset, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	lr := &io.LimitedReader{R: f, N: part.Size}
	// stats for this part
	stats := make(map[string]Stats)

	scanner := bufio.NewScanner(lr)
	for scanner.Scan() {
		line := scanner.Text()
		name, tempStr, found := strings.Cut(line, ";")
		if !found {
			continue
		}
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		s, ok := stats[name]
		if !ok {
			s.Min = temp
			s.Max = temp
			s.Sum = temp
			s.Count = 1
		} else {
			s.Min = min(s.Min, temp)
			s.Max = max(s.Max, temp)
			s.Sum += temp
			s.Count++
		}
		stats[name] = s // note: XXX: double hashing
	}
	resultC <- stats
}

// splitFile partitions a file into a fixed number of chunks, breaking on lines.
func splitFile(fn string, numParts int) ([]Part, error) {
	const maxLineLength = 100
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Get file size.
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var (
		size      = st.Size()
		splitSize = size / int64(numParts)
		buf       = make([]byte, maxLineLength)
		parts     = make([]Part, 0, numParts)
		offset    = int64(0)
	)
	for i := 0; i < numParts; i++ {
		if i == numParts-1 {
			// This is the last part, if we have remaining data [...]
			if offset < size {
				parts = append(parts, Part{Offset: offset, Size: size - offset})
			}
			break
		}
		seekOffset := max(offset+splitSize-maxLineLength, 0)
		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(f, buf)
		chunk := buf[:n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			return nil, fmt.Errorf("newline not found")
		}
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, Part{Offset: offset, Size: nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
}

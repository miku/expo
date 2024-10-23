// Partition the file and use goroutines to work on each part.
//
// real    1m21.538s
// user    9m43.259s
// sys     0m21.870s

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
	"runtime"
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

type Part struct {
	Offset int64
	Length int64
}

// Measurements, as there is no need to keep all numbers around, we can compute
// them on the fly.
type Measurements struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

// Add adds a measurement.
func (m *Measurements) Add(v float64) {
	if v > m.Max {
		m.Max = v
	} else if v < m.Min {
		m.Min = v
	}
	m.Sum = m.Sum + v
	m.Count++
}

// Merge merges another measurement into the current one.
func (m *Measurements) Merge(o *Measurements) {
	if o.Min < m.Min {
		m.Min = o.Min
	}
	if o.Max > m.Max {
		m.Max = o.Max
	}
	m.Sum = m.Sum + o.Sum
	m.Count = m.Count + o.Count
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
	parts, err := partitionFile(*filename, runtime.NumCPU())
	if err != nil {
		log.Fatal(err)
	}
	var resultC = make(chan map[string]*Measurements)
	for _, part := range parts {
		go readPart(*filename, part, resultC)
	}
	data := make(map[string]*Measurements)
	for i := 0; i < len(parts); i++ {
		result := <-resultC
		for k, v := range result {
			if data[k] == nil {
				data[k] = v
			} else {
				data[k].Merge(v)
			}
		}
	}
	// At this point, data contains the merged data from all measurements.
	keys := slices.Collect(maps.Keys(data))
	sort.Strings(keys)
	for _, k := range keys {
		avg := data[k].Sum / float64(data[k].Count)
		fmt.Printf("%s\t%0.2f/%0.2f/%0.2f\n", k, data[k].Min, avg, data[k].Max)
	}
}

// partitionFile returns a slice of offset, length pairs that split the file up
// at newline boundaries.
func partitionFile(fn string, numPartitions int) ([]Part, error) {
	const maxLineLength = 100
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var (
		size      = st.Size()
		splitSize = size / int64(numPartitions)
		buf       = make([]byte, maxLineLength)
		parts     = make([]Part, 0, numPartitions)
		offset    = int64(0)
	)
	for i := 0; i < numPartitions; i++ {
		if i == numPartitions-1 {
			// This is the last part, if we have remaining data [...]
			if offset < size {
				parts = append(parts, Part{Offset: offset, Length: size - offset})
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
		parts = append(parts, Part{Offset: offset, Length: nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
}

// readPart works on a subset of the lines in a file.
func readPart(fn string, part Part, resultC chan map[string]*Measurements) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Seek(part.Offset, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	var data = make(map[string]*Measurements)
	lr := io.LimitReader(f, part.Length)
	scanner := bufio.NewScanner(lr)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			log.Fatalf("expected two fields: %s, got %d", line, len(parts))
		}
		name := strings.TrimSpace(parts[0])
		temp, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			log.Fatalf("invalid temp: %s", parts[1])
		}
		if _, ok := data[name]; !ok {
			data[name] = &Measurements{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		} else {
			data[name].Add(temp)
		}
	}
	resultC <- data
}

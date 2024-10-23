// Do not allocate a whole slice for keeping all measurements in memory.
//
// i7-8550U
//
// real    3m57.385s
// user    3m50.386s
// sys     0m9.299s

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

// Measurements per station.
type Measurements struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func (m *Measurements) Add(v float64) {
	if v > m.Max {
		m.Max = v
	} else if v < m.Min {
		m.Min = v
	}
	m.Sum = m.Sum + v
	m.Count++
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
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var data = make(map[string]*Measurements)
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
	keys := maps.Keys(data)
	sort.Strings(keys)
	for _, k := range keys {
		avg := data[k].Sum / float64(data[k].Count)
		fmt.Printf("%s\t%0.2f/%0.2f/%0.2f\n", k, data[k].Min, data[k].Max, avg)
	}
}

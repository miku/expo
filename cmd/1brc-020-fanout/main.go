// Basic parallel version, fan-out-fan-in pattern. Not using other optimizations.
//
// real    2m39.037s
// user    10m2.609s
// sys     0m28.006s

package main

import (
	"bufio"
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
	"sync"

	"golang.org/x/exp/maps"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

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

func worker(queue chan []string, result chan map[string]*Measurements, wg *sync.WaitGroup) {
	defer wg.Done()
	var data = make(map[string]*Measurements)
	for batch := range queue {
		for _, line := range batch {
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
	}
	result <- data
}

func merger(data map[string]*Measurements, result chan map[string]*Measurements, done chan bool) {
	for m := range result {
		for k, v := range m {
			if _, ok := data[k]; !ok {
				data[k] = &Measurements{
					Min:   v.Min,
					Max:   v.Max,
					Sum:   v.Sum,
					Count: v.Count,
				}
			} else {
				data[k].Merge(v)
			}
		}
	}
	done <- true
}

func main() {
	flag.Parse()
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	var (
		batchSize = 20_000_000
		queue     = make(chan []string)
		result    = make(chan map[string]*Measurements)
		wg        sync.WaitGroup
		done      = make(chan bool)
		// accumulate all results here
		data = make(map[string]*Measurements)
	)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(queue, result, &wg)
	}
	go merger(data, result, done)
	// start reading the file and fan out
	br := bufio.NewReader(f)
	batch := make([]string, 0)
	i := 0
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		batch = append(batch, line)
		i++
		if i%batchSize == 0 {
			b := make([]string, len(batch))
			copy(b, batch)
			queue <- b
			batch = nil
		}
	}
	queue <- batch // rest, no copy required
	close(queue)
	wg.Wait()
	close(result)
	<-done
	// At this point, data contains the merged data from all measurements.
	keys := maps.Keys(data)
	sort.Strings(keys)
	for _, k := range keys {
		avg := data[k].Sum / float64(data[k].Count)
		fmt.Printf("%s\t%0.2f/%0.2f/%0.2f\n", k, data[k].Min, avg, data[k].Max)
	}
}

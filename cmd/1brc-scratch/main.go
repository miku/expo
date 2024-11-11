// Playground to implement various ideas.
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

func main() {
	flag.Parse()
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	data := make(map[string]*Measurements)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimRight(line, "\n")
		parts := strings.Split(line, ";")
		temp, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		name := parts[0]
		if _, ok := data[name]; ok {
			data[name].Add(temp)
		} else {
			data[name] = &Measurements{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		}

	}

}

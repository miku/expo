// Baseline with an unoptimized scanner, limited memory; avoid string.Split.
//
// i7-8550U
//
// real    2m17.561s
// user    2m10.390s
// sys     0m7.209s

package main

import (
	"bufio"
	"flag"
	"fmt"
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

// Station data.
type Station struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func ToResultString(m map[string]*Station) string {
	var (
		sb   strings.Builder
		keys = slices.Collect(maps.Keys(m))
	)
	_, _ = sb.WriteString("{")
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		s := fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", k, v.Min, v.Sum/float64(v.Count), v.Max)
		_, _ = sb.WriteString(s)
	}
	_, _ = sb.WriteString("}")
	return sb.String()
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
	var data = map[string]*Station{}
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var (
			line  = scanner.Text()
			index = strings.Index(line, ";")
			name  = strings.TrimSpace(line[:index])
			tempS = strings.TrimSpace(line[index+1:])
		)
		temp, err := strconv.ParseFloat(tempS, 64)
		if err != nil {
			log.Fatal(err)
		}
		st, ok := data[name]
		switch {
		case !ok:
			data[name] = &Station{
				Name:  name,
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		default:
			if temp < st.Min {
				st.Min = temp
			}
			if temp > st.Max {
				st.Max = temp
			}
			st.Sum += temp
			st.Count++
		}
	}
	s := ToResultString(data)
	fmt.Println(s)
}

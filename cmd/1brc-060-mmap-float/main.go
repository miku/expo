// With mmap, but use a tweak for parsing a float; taken from https://github.com/valyala/fastjson/blob/6dae91c8e11a7fa6a257a550b75cba53ab81693e/fastfloat/parse.go#L203
//
// real    0m41.969s
// user    4m35.329s
// sys     0m11.887s

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/mmap"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

const chunkSize = 67108864

// Exact powers of 10.
//
// This works faster than math.Pow10, since it avoids additional multiplication.
var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16,
}

// Measurements, as there is no need to keep all numbers around, we can compute
// them on the fly.
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

// aggregate aggregates measurements by reading a chunk from an io.ReaderAt and
// passing the result to a results channel.
func aggregate(rat io.ReaderAt, offset, length int, resultC chan map[string]*Measurements, sem chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if length == 0 {
		return
	}
	buf := make([]byte, length)
	_, err := rat.ReadAt(buf, int64(offset))
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println(offset, length)
	var (
		data    = make(map[string]*Measurements)
		j, k, l = 0, 0, 0 // j=start, k=semi, l=newline
		n       = 0
	)
	for i := 0; i < length; i++ {
		if buf[i] == ';' {
			k = i
		} else if buf[i] == '\n' {
			l = i
			// TODO: get rid of strings!
			name := string(buf[j:k])
			temp := ParseBestEffort(string(buf[k+1 : l]))
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
			j = l + 1
			n++
		}
	}
	resultC <- data
	<-sem
}

func merger(data map[string]*Measurements, resultC chan map[string]*Measurements, done chan bool) {
	for m := range resultC {
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
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	r, err := mmap.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	var (
		resultC = make(chan map[string]*Measurements)
		done    = make(chan bool)
		sem     = make(chan bool, runtime.NumCPU())
		wg      sync.WaitGroup
		data    = make(map[string]*Measurements)
	)
	go merger(data, resultC, done)
	var i, j int // start and stop index
	for i < r.Len() {
		j = i + chunkSize
		if j > r.Len() {
			L := j - i
			wg.Add(1)
			sem <- true
			go aggregate(r, i, L, resultC, sem, &wg)
			break
		}
		for {
			// TODO: maybe split this into goroutines as well
			if r.At(j) == '\n' {
				break // found newline
			}
			j++
		}
		L := j - i
		wg.Add(1)
		sem <- true
		go aggregate(r, i, L, resultC, sem, &wg)
		i = j + 1
	}
	wg.Wait()
	close(resultC)
	<-done
	keys := maps.Keys(data)
	sort.Strings(keys)
	for _, k := range keys {
		avg := data[k].Sum / float64(data[k].Count)
		fmt.Printf("%s\t%0.2f/%0.2f/%0.2f\n", k, data[k].Min, data[k].Max, avg)
	}
}

// ParseBestEffort parses floating-point number s.
//
// It is equivalent to strconv.ParseFloat(s, 64), but is faster.
//
// 0 is returned if the number cannot be parsed.
// See also Parse, which returns parse error if the number cannot be parsed.
func ParseBestEffort(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	i := uint(0)
	minus := s[0] == '-'
	if minus {
		i++
		if i >= uint(len(s)) {
			return 0
		}
	}

	d := uint64(0)
	j := i
	for i < uint(len(s)) {
		if s[i] >= '0' && s[i] <= '9' {
			d = d*10 + uint64(s[i]-'0')
			i++
			if i > 18 {
				// The integer part may be out of range for uint64.
				// Fall back to slow parsing.
				f, err := strconv.ParseFloat(s, 64)
				if err != nil && !math.IsInf(f, 0) {
					return 0
				}
				return f
			}
			continue
		}
		break
	}
	if i <= j {
		s = s[i:]
		if strings.HasPrefix(s, "+") {
			s = s[1:]
		}
		// "infinity" is needed for OpenMetrics support.
		// See https://github.com/OpenObservability/OpenMetrics/blob/master/OpenMetrics.md
		// if strings.EqualFold(s, "inf") || strings.EqualFold(s, "infinity") {
		// 	if minus {
		// 		return -inf
		// 	}
		// 	return inf
		// }
		// if strings.EqualFold(s, "nan") {
		// 	return nan
		// }
		return 0
	}
	f := float64(d)
	if i >= uint(len(s)) {
		// Fast path - just integer.
		if minus {
			f = -f
		}
		return f
	}

	if s[i] == '.' {
		// Parse fractional part.
		i++
		if i >= uint(len(s)) {
			return 0
		}
		k := i
		for i < uint(len(s)) {
			if s[i] >= '0' && s[i] <= '9' {
				d = d*10 + uint64(s[i]-'0')
				i++
				if i-j >= uint(len(float64pow10)) {
					// The mantissa is out of range. Fall back to standard parsing.
					f, err := strconv.ParseFloat(s, 64)
					if err != nil && !math.IsInf(f, 0) {
						return 0
					}
					return f
				}
				continue
			}
			break
		}
		if i < k {
			return 0
		}
		// Convert the entire mantissa to a float at once to avoid rounding errors.
		f = float64(d) / float64pow10[i-k]
		if i >= uint(len(s)) {
			// Fast path - parsed fractional number.
			if minus {
				f = -f
			}
			return f
		}
	}
	if s[i] == 'e' || s[i] == 'E' {
		// Parse exponent part.
		i++
		if i >= uint(len(s)) {
			return 0
		}
		expMinus := false
		if s[i] == '+' || s[i] == '-' {
			expMinus = s[i] == '-'
			i++
			if i >= uint(len(s)) {
				return 0
			}
		}
		exp := int16(0)
		j := i
		for i < uint(len(s)) {
			if s[i] >= '0' && s[i] <= '9' {
				exp = exp*10 + int16(s[i]-'0')
				i++
				if exp > 300 {
					// The exponent may be too big for float64.
					// Fall back to standard parsing.
					f, err := strconv.ParseFloat(s, 64)
					if err != nil && !math.IsInf(f, 0) {
						return 0
					}
					return f
				}
				continue
			}
			break
		}
		if i <= j {
			return 0
		}
		if expMinus {
			exp = -exp
		}
		f *= math.Pow10(int(exp))
		if i >= uint(len(s)) {
			if minus {
				f = -f
			}
			return f
		}
	}
	return 0
}

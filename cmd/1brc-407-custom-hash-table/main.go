// R7: custom hash map
//
// # Intel i7-8550U
//
// real    0m28.600s
// user    0m23.561s
// sys     0m3.576s
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sort"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

type Stats struct {
	Min, Max, Count int32
	Sum             int64
}

type Item struct {
	key   []byte
	stats *Stats
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
	if err := R7(*filename, bw); err != nil {
		log.Fatal(err)
	}
}

func R7(fn string, w io.Writer) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	const numBuckets = 1 << 17
	items := make([]Item, numBuckets)
	size := 0
	// read buffer
	buf := make([]byte, 1<<20)
	readStart := 0
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
			const (
				offset64 = 14695981039346656037
				prime64  = 1099511628211
			)
			var name, after []byte
			hash := uint64(offset64)
			var i int
			for ; i < len(chunk); i++ {
				c := chunk[i]
				if c == ';' {
					name = chunk[:i]
					after = chunk[i+1:]
					break
				}
				hash ^= uint64(c) // FNV1a is XOR then
				hash *= prime64
			}
			if i == len(chunk) {
				break
			}

			idx := 0
			neg := false
			if after[idx] == '-' {
				neg = true
				i++
			}
			temp := int32(after[idx] - '0')
			idx++
			if after[idx] != '.' {
				temp = temp*10 + int32(after[idx]-'0')
				idx++
			}
			idx++ // skip '.'
			temp = temp*10 + int32(after[idx]-'0')
			idx += 2 // skip last digit and '\n'
			if neg {
				temp = -temp
			}
			chunk = after[idx:]

			hashIndex := int(hash & uint64(numBuckets-1))
			for {
				if items[hashIndex].key == nil {
					// empty slot, add new item
					key := make([]byte, len(name))
					copy(key, name)
					items[hashIndex] = Item{
						key: key,
						stats: &Stats{
							Min:   temp,
							Max:   temp,
							Sum:   int64(temp),
							Count: 1,
						},
					}
					size++
					if size > numBuckets/2 {
						panic("too many items in hash table")
					}
					break
				}
				if bytes.Equal(items[hashIndex].key, name) {
					// found matching slot
					s := items[hashIndex].stats
					s.Min = min(s.Min, temp)
					s.Max = max(s.Max, temp)
					s.Sum += int64(temp)
					s.Count++
					break
				}
				hashIndex++
				if hashIndex > numBuckets {
					hashIndex = 0
				}
			}
		}
		readStart = copy(buf, remaining)
	}
	stationItems := make([]Item, 0, size)
	for _, item := range items {
		if item.key == nil {
			continue
		}
		stationItems = append(stationItems, item)
	}
	sort.Slice(stationItems, func(i, j int) bool {
		return string(stationItems[i].key) < string(stationItems[j].key)
	})

	fmt.Fprint(w, "{")
	for i, item := range stationItems {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		s := item.stats
		mean := float64(s.Sum) / float64(s.Count) / 10
		fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f", item.key, float64(s.Min)/10, mean, float64(s.Max)/10)
	}
	fmt.Fprint(w, "}\n")
	return nil
}

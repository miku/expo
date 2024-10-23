// R10: fast semicolon find with SWAR, plus other optimizations
//
// CPU: Intel i7-8550U
//
// real    0m7.961s
// user    0m40.479s
// sys     0m6.662s
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

const BroadcastSemicolon = 0x3B3B3B3B3B3B3B3B
const Broadcast0x01 = 0x0101010101010101
const Broadcast0x80 = 0x8080808080808080

type Stats struct {
	Min, Max, Count int32
	Sum             int64
}

type Item struct {
	Key   []byte
	Stats *Stats
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
	if err := R10(*filename, bw); err != nil {
		log.Fatal(err)
	}
}

func R10(fn string, w io.Writer) error {
	parts, err := splitFile(fn, runtime.NumCPU())
	if err != nil {
		return err
	}
	// run one goroutine per part
	resultC := make(chan map[string]*Stats)
	for _, part := range parts {
		go processPart(fn, part, resultC) // XXX: cancellation and wait?
	}

	// merge stats from goroutines in totals
	totals := make(map[string]*Stats)
	for i := 0; i < len(parts); i++ {
		result := <-resultC
		for name, s := range result {
			ts := totals[name]
			if ts == nil {
				totals[name] = s
				continue
			}
			ts.Min = min(ts.Min, s.Min)
			ts.Max = max(ts.Max, s.Max)
			ts.Sum += s.Sum
			ts.Count += s.Count
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
		mean := float64(s.Sum) / float64(s.Count) / 10
		fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f", name, float64(s.Min)/10, mean, float64(s.Max)/10)
	}
	fmt.Fprint(w, "}\n")
	return nil
}

// processPart does not do any error handling for the moment; optimized.
func processPart(fn string, part Part, resultC chan map[string]*Stats) {
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

	const (
		// use custom hash map
		numBuckets = 1 << 17
		// FNV-1 64-bit constants from hash/fnv.
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	var (
		items     = make([]Item, numBuckets)
		size      = 0
		buf       = make([]byte, 1<<20) // 1MB
		readStart = 0
	)
	for {
		n, err := lr.Read(buf[readStart:])
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if readStart+n == 0 {
			break // nothing to read
		}
		chunk := buf[:readStart+n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			break // no more newlines
		}
		remaining := chunk[newline+1:]
		chunk = chunk[:newline+1]

	chunkLoop:
		for {
			// Faster semicolon hashing.
			var hash uint64
			var name, after []byte
			if len(chunk) < 8 {
				break chunkLoop
			}
			// Turns eight bytes into a uint64
			nameWord0 := binary.NativeEndian.Uint64(chunk)
			matchBits := semicolonMatchBits(nameWord0)
			if matchBits != 0 {
				// semicolon is in the first 8 bytes
				nameLen := calcNameLen(matchBits)
				nameWord0 = maskWord(nameWord0, matchBits)
				name = chunk[:nameLen]
				after = chunk[nameLen+1:]
				hash = calcHash(nameWord0)
			} else {
				// station name is longer so keep looking for the semicolon in
				// uint64 chunks
				nameLen := 8
				hash = calcHash(nameWord0)
				for {
					if nameLen > len(chunk)-8 {
						break chunkLoop
					}
					lastNameWord := binary.NativeEndian.Uint64(chunk[nameLen:])
					matchBits = semicolonMatchBits(lastNameWord)
					if matchBits != 0 {
						nameLen += calcNameLen(matchBits)
						name = chunk[:nameLen]
						after = chunk[nameLen+1:]
						break
					}
					nameLen += 8
				}
			}
			index := 0
			negative := false
			if after[index] == '-' {
				negative = true
				index++
			}
			temp := int32(after[index] - '0')
			index++
			if after[index] != '.' {
				temp = temp*10 + int32(after[index]-'0')
				index++
			}
			index++ // skip '.'
			temp = temp*10 + int32(after[index]-'0')
			index += 2 // skip last digit and '\n'
			if negative {
				temp = -temp
			}
			chunk = after[index:]
			hashIndex := int(hash & (numBuckets - 1))
			for {
				if items[hashIndex].Key == nil {
					// Found empty slot, add new item (copying key).
					key := make([]byte, len(name))
					copy(key, name)
					items[hashIndex] = Item{
						Key: key,
						Stats: &Stats{
							Min:   temp,
							Max:   temp,
							Sum:   int64(temp),
							Count: 1,
						},
					}
					size++
					if size > numBuckets/2 {
						log.Fatal("too many items in hash table")
					}
					break
				}
				if bytes.Equal(items[hashIndex].Key, name) {
					// Found matching slot, add to existing stats.
					s := items[hashIndex].Stats
					s.Min = min(s.Min, temp)
					s.Max = max(s.Max, temp)
					s.Sum += int64(temp)
					s.Count++
					break
				}
				// Slot already holds another key, try next slot (linear probe).
				hashIndex++
				if hashIndex >= numBuckets {
					hashIndex = 0
				}
			}
		}
		readStart = copy(buf, remaining)
	}
	result := make(map[string]*Stats, size)
	for _, item := range items {
		if item.Key == nil {
			continue
		}
		result[string(item.Key)] = item.Stats
	}
	resultC <- result
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

func calcNameLen(b uint64) int {
	return (bits.TrailingZeros64(b) >> 3)
}

func calcHash(word uint64) uint64 {
	return bits.RotateLeft64(word*0x51_7c_c1_b7_27_22_0a_95, 17)
}

// semicolonMatchBits locates the ASCII semicolon byte and returns a long with bits set to one where it was found.
func semicolonMatchBits(word uint64) uint64 {
	diff := word ^ BroadcastSemicolon
	return (diff - Broadcast0x01) & (^diff & Broadcast0x80)
}

func maskWord(word, matchBits uint64) uint64 {
	mask := matchBits ^ (matchBits - 1)
	return word & mask
}

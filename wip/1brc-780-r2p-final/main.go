package main

import (
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/dolthub/swiss"
)

const READ_BUFFER_SIZE = 2048 * 2048
const N_WORKERS = 75

type TrashItem struct {
	Idx     int
	Value   []byte
	Initial bool
}

type StationData struct {
	Name  string
	Min   int
	Max   int
	Sum   int
	Count int
}

var lock = &sync.Mutex{}
var lockIdx = 0

func trashBin(input chan *TrashItem, output chan *swiss.Map[uint64, *StationData], wg *sync.WaitGroup) {
	defer wg.Done()
	data := swiss.NewMap[uint64, *StationData](1024)

	can := []*TrashItem{}
	buffer := make([]byte, 1024)

	for item := range input {
		can = append(can, item)
		can = saveCan(can, data, buffer)
	}

	output <- data
}

func saveCan(can []*TrashItem, data *swiss.Map[uint64, *StationData], buffer []byte) []*TrashItem {
	for i, ref := range can {
		if ref.Idx == 0 {
			_, nameInit, nameEnd, tempInit, tempEnd := nextLine(0, ref.Value)
			processLine(ref.Value[nameInit:nameEnd], ref.Value[tempInit:tempEnd], data)
			return slices.Delete(can, i, i+1)
		}

		for j, oth := range can {
			if ref.Idx == oth.Idx && i != j {
				if ref.Initial {
					copy(buffer[:len(ref.Value)], ref.Value)
					copy(buffer[len(ref.Value):], oth.Value)
				} else {
					copy(buffer[:len(oth.Value)], oth.Value)
					copy(buffer[len(oth.Value):], ref.Value)
				}
				total := len(ref.Value) + len(oth.Value)

				end, nameInit, nameEnd, tempInit, tempEnd := nextLine(0, buffer)
				processLine(buffer[nameInit:nameEnd], buffer[tempInit:tempEnd], data)

				if end < total {
					_, nameInit, nameEnd, tempInit, tempEnd := nextLine(end, buffer)
					processLine(buffer[nameInit:nameEnd], buffer[tempInit:tempEnd], data)
				}

				if i > j {
					can = slices.Delete(can, i, i+1)
					can = slices.Delete(can, j, j+1)
				} else {
					can = slices.Delete(can, j, j+1)
					can = slices.Delete(can, i, i+1)
				}

				return can
			}
		}
	}

	return can
}

func consumer(file *os.File, trash chan *TrashItem, output chan *swiss.Map[uint64, *StationData], wg *sync.WaitGroup) {
	defer wg.Done()
	data := swiss.NewMap[uint64, *StationData](1024)

	readBuffer := make([]byte, READ_BUFFER_SIZE)
	for {
		lock.Lock()
		lockIdx++
		idx := lockIdx
		n, err := file.Read(readBuffer)
		lock.Unlock()

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// ignoring first line
		start := 0
		for i := 0; i < n; i++ {
			if readBuffer[i] == 10 {
				start = i + 1
				break
			}
		}
		trash <- &TrashItem{idx - 1, readBuffer[:start], false}

		// ignoring last line
		final := 0
		for i := n - 1; i >= 0; i-- {
			if readBuffer[i] == 10 {
				final = i
				break
			}
		}
		trash <- &TrashItem{idx, readBuffer[final+1 : n], true}

		readingIndex := start
		for readingIndex < final {
			next, nameInit, nameEnd, tempInit, tempEnd := nextLine(readingIndex, readBuffer)
			readingIndex = next
			processLine(readBuffer[nameInit:nameEnd], readBuffer[tempInit:tempEnd], data)
		}
	}

	output <- data
}

func nextLine(readingIndex int, reading []byte) (nexReadingIndex, nameInit, nameEnd, tempInit, tempEnd int) {
	i := readingIndex
	nameInit = readingIndex
	for reading[i] != 59 { // ;
		i++
	}
	nameEnd = i

	i++ // skip ;

	tempInit = i
	for i < len(reading) && reading[i] != 10 { // \n
		i++
	}
	tempEnd = i

	readingIndex = i + 1
	return readingIndex, nameInit, nameEnd, tempInit, tempEnd
}

func processLine(name, temperature []byte, data *swiss.Map[uint64, *StationData]) {
	temp := bytesToInt(temperature)
	id := hash(name)
	station, ok := data.Get(id)
	if !ok {
		data.Put(id, &StationData{string(name), temp, temp, temp, 1})
	} else {
		if temp < station.Min {
			station.Min = temp
		}
		if temp > station.Max {
			station.Max = temp
		}
		station.Sum += temp
		station.Count++
	}
}

func run() {
	outputChannels := make([]chan *swiss.Map[uint64, *StationData], N_WORKERS+1)

	// Read file
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	var wgTrash sync.WaitGroup

	wg.Add(N_WORKERS)
	wgTrash.Add(1)
	trash := make(chan *TrashItem, N_WORKERS*2)
	output := make(chan *swiss.Map[uint64, *StationData], 1)
	go trashBin(trash, output, &wgTrash)
	outputChannels[0] = output

	for i := 0; i < N_WORKERS; i++ {
		output := make(chan *swiss.Map[uint64, *StationData], 1)
		go consumer(file, trash, output, &wg)
		outputChannels[i+1] = output
	}

	wg.Wait()
	close(trash)
	wgTrash.Wait()

	for i := 0; i < N_WORKERS+1; i++ {
		close(outputChannels[i])
	}

	data := swiss.NewMap[uint64, *StationData](1000)
	for i := 0; i < N_WORKERS+1; i++ {
		m := <-outputChannels[i]
		m.Iter(func(station uint64, stationData *StationData) bool {
			v, ok := data.Get(station)
			if !ok {
				data.Put(station, stationData)
			} else {
				if stationData.Min < v.Min {
					v.Min = stationData.Min
				}
				if stationData.Max > v.Max {
					v.Max = stationData.Max
				}
				v.Sum += stationData.Sum
				v.Count += stationData.Count
			}

			return false
		})
	}

	printResult(data)
}

func hash(name []byte) uint64 {
	var h uint64 = 5381
	for _, b := range name {
		h = (h << 5) + h + uint64(b)
	}
	return h
}

func printResult(data *swiss.Map[uint64, *StationData]) {
	result := make(map[string]*StationData, data.Count())
	keys := make([]string, 0, data.Count())

	data.Iter(func(k uint64, v *StationData) (stop bool) {
		keys = append(keys, v.Name)
		result[v.Name] = v
		return false
	})
	sort.Strings(keys)

	print("{")
	for _, k := range keys {
		v := result[k]
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", k, float64(v.Min)/10, (float64(v.Sum)/10)/float64(v.Count), float64(v.Max)/10)
	}
	print("}\n")
}

func bytesToInt(byteArray []byte) int {
	var result int
	negative := false

	for _, b := range byteArray {
		if b == 46 { // .
			continue
		}

		if b == 45 { // -
			negative = true
			continue
		}
		result = result*10 + int(b-48)
	}

	if negative {
		return -result
	}

	return result
}

func main() {
	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	started := time.Now()
	run()
	fmt.Printf("%0.6f\n", time.Since(started).Seconds())
}

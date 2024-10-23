package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"
)

var size = flag.Int("b", 1024, "buffer size")

type opts struct {
	size int
}

func consumer(ch chan []byte) {
	for {
		<-ch
	}
}

func run(opts *opts) error {
	buf := make([]byte, opts.size)
	leftoverBuffer := make([]byte, 1024) // TODO: magic number
	leftoverSize := 0
	file, err := os.Open("measurements.txt")
	if err != nil {
		return err
	}
	ch := make(chan []byte, 10)
	go consumer(ch)
	defer file.Close()
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		m := 0
		for i := n - 1; i >= 0; i-- {
			if buf[i] == 10 {
				m = i
				break
			}
		}
		data := make([]byte, m+leftoverSize)
		copy(data, leftoverBuffer[:leftoverSize])
		copy(data[leftoverSize:], buf[:m])
		copy(leftoverBuffer, buf[m+1:n])
		leftoverSize = n - m - 1
		ch <- data
	}
	return nil
}

func main() {
	flag.Parse()
	opts := &opts{size: *size}
	started := time.Now()
	if err := run(opts); err != nil {
		log.Fatal(err)
	}
	log.Printf("%0.6f", time.Since(started).Seconds())
}

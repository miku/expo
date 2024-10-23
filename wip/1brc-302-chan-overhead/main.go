package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"
)

// $ ./1brc-301-file-read -b 1048576 && ./1brc-302-chan-overhead -b 1048576 && ./1brc-303-chan-copy -b 1048576
// 2024/11/05 00:31:11 1.572377
// 2024/11/05 00:31:12 1.587595
// 2024/11/05 00:31:18 5.359889

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
	ch := make(chan []byte, 10)
	go consumer(ch)
	file, err := os.Open("measurements.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, opts.size)
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ch <- buf
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

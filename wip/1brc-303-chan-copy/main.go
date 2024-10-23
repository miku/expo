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
	ch := make(chan []byte, 10)
	go consumer(ch)
	file, err := os.Open("measurements.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, opts.size)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		data := make([]byte, n)
		copy(data, buf[:n])
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

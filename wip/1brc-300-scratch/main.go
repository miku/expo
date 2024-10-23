package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"
)

var size = flag.Int("b", 1024, "buffer size")

type opts struct {
	size int
}

func run(opts *opts) error {
	file, err := os.Open("measurements.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, opts.size), opts.size)
	for scanner.Scan() {
		scanner.Bytes()
	}
	if err := scanner.Err(); err != nil {
		return err
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

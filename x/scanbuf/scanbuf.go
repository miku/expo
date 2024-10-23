package scanbuf

import (
	"bufio"
	"log"
	"os"
)

// WithBuffer will read data from a file with a buffer of size.
func WithBuffer(filename string, size int) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	buf := make([]byte, size)
	scanner.Buffer(buf, size)
	for scanner.Scan() {
		scanner.Bytes()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

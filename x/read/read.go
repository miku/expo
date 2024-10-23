package read

import (
	"bufio"
	"io"
	"os"
)

// Raw reads from a file directly, using a buffer of a given size.
func Raw(filename string, size int) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for {
		var p []byte = make([]byte, size)
		if _, err := f.Read(p); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
	return nil
}

// ReadByte uses bufio Read, reading a single byte at a time.
func ReadByte(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		_, err := br.ReadByte()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// ReadString uses buffered io and reads a line into a string.
func ReadString(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		_, err := br.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// ReadString uses buffered io and reads a line into a byte slice.
func ReadBytes(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		_, err := br.ReadBytes('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

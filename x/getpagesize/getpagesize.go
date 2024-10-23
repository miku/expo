package main

import (
	"fmt"
	"os"
)

func main() {
	pageSize := os.Getpagesize()
	fmt.Printf("Page size: %d bytes (%d KB)\n", pageSize, pageSize/1024)
}

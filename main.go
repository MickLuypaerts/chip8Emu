package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
}

func usage() {
	fmt.Printf("Usage: %s [OPTION] [PATTERN] [FILE]\n", os.Args[0])
	fmt.Printf("Use %s -help for a list of flags.\n", os.Args[0])
}

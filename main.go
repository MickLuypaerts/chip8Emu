package main

import (
	"chip8/chip8"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	var chip8 chip8.Chip8
}

func usage() {
	fmt.Printf("Usage: %s [OPTION] [PATTERN] [FILE]\n", os.Args[0])
	fmt.Printf("Use %s -help for a list of flags.\n", os.Args[0])
}

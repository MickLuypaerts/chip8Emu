package main

import (
	"chip8/chip8"
	"chip8/emulator"
	"chip8/view"
	"fmt"
	"log"
	"os"
)

func main() {
	chip := new(chip8.Chip8)
	tui := new(view.TUI)
	emu, err := emulator.CreateEmulator(os.Args, "q", chip, tui)
	if err != nil {
		log.Fatal(err)
	}
	emu.Run()

}
func usage() {
	fmt.Printf("Usage: %s [FILE]\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("Emulator Controls:\n")
	fmt.Printf("|key|    function   |\n")
	fmt.Printf("|---|---------------|\n")
	fmt.Printf("| q |      quit     |\n")
	fmt.Printf("| s |    1 cycle    |\n")
	fmt.Printf("| r |  run program  |\n")
	fmt.Printf("| R | stop program  |\n")
	fmt.Printf("| j | Mem map down  |\n")
	fmt.Printf("| k |  Mem map up   |\n")
	fmt.Printf("| gg|  Mem map top  |\n")
	fmt.Printf("| G |Mem map bottom |\n")
}

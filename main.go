package main

import (
	"chip8/chip8"
	"chip8/emulator"
	"chip8/view"
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

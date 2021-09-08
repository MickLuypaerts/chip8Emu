package main

import (
	"log"
	"os"

	"github.com/MickLuypaerts/chip8Emu/chip8"
	"github.com/MickLuypaerts/chip8Emu/emulator"
	"github.com/MickLuypaerts/chip8Emu/view"
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

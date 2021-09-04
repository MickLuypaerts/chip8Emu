package main

import (
	"chip8/chip8"
	"chip8/view"
	"fmt"
	"log"
	"os"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	var chip8 chip8.Chip8
	var TUI view.TUI
	err := chip8.Init(os.Args[1], TUI.UpdateScreen, TUI.SetEmuInfo, TUI.SetKeyInfo)
	if err != nil {
		log.Fatalf("failed to initialize chip8: %v", err)
	}

	// TUI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	TUI.Init(&chip8)
	ui.Render(TUI.Grid)
	previousKey := ""
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":

			return
		case "j", "<Down>":
			view.ScrollDown(TUI.LMem)

		case "k", "<Up>":
			view.ScrollUp(TUI.LMem)
		case "g":
			if previousKey == "g" {
				view.ScrollTop(TUI.LMem)
				previousKey = ""
			}
		case "G", "<End>":
			view.ScrollBottom(TUI.LMem)
		case "s":
			chip8.EmulateCycle()
		case "r":
			chip8.Run()
		case "R":
			chip8.Stop()
		case "1":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x1)
		case "2":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x2)
		case "3":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x3)
		case "4":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x4)
		case "5":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x5)
		case "6":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x6)
		case "7":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x7)
		case "8":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x8)
		case "9":
			sendKeyboardInterrupt(chip8.KeyBoardInterrupt, 0x9)
		}
		previousKey = e.ID
	}

}
func sendKeyboardInterrupt(c chan byte, key byte) {
	c <- key
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

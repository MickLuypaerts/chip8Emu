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
	err := chip8.Init(os.Args[1])
	if err != nil {
		log.Fatalf("failed to initialize chip8: %v", err)
	}

	// TUI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var TUI view.TUI
	TUI.Init(&chip8)
	ui.Render(TUI.Grid)
	previousKey := ""
	uiEvents := ui.PollEvents()
	chip8.ScreenFunc = TUI.UpdateScreen
	chip8.InfoFunc = TUI.SetEmuInfo
	chip8.KeyFunc = TUI.SetKeyInfo

	for {
		chip8.ClearKeys()
		e := <-uiEvents

		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			view.ScrollDown(TUI.LMem)
			chip8.ClearKeys()
		case "k", "<Up>":
			view.ScrollUp(TUI.LMem)
			chip8.ClearKeys()
		case "g":
			if previousKey == "g" {
				view.ScrollTop(TUI.LMem)
			}
		case "G", "<End>":
			view.ScrollBottom(TUI.LMem)
			chip8.ClearKeys()
		case "s":
			chip8.EmulateCycle()
			chip8.ClearKeys()
		case "r":
			chip8.Run()
			chip8.ClearKeys()
		case "R":
			chip8.Stop()
			chip8.ClearKeys()
		case "1":
			chip8.PressKey(0x1)
		case "2":
			chip8.PressKey(0x2)
		case "3":
			chip8.PressKey(0x3)
		case "4":
			chip8.PressKey(0x4)
		case "5":
			chip8.PressKey(0x5)
		case "6":
			chip8.PressKey(0x6)
		case "7":
			chip8.PressKey(0x7)
		case "8":
			chip8.PressKey(0x8)
		case "9":
			chip8.PressKey(0x9)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

	}

}

func usage() {
	fmt.Printf("Usage: %s [FILE]\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("Controls:\n")
	fmt.Printf("|key|  function  |\n")
	fmt.Printf("|---|------------|\n")
	fmt.Printf("| q |    quit    |\n")
	fmt.Printf("| s |  1 cycle   |\n")
	fmt.Printf("| r |run program |\n")
	fmt.Printf("| R |stop program|\n")
	fmt.Printf("| j |Mem map down|\n")
	fmt.Printf("| k | Mem map up |\n")
	fmt.Printf("| gg|Mem map top |\n")
	fmt.Printf("| G |Mem map top |\n")
}

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
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			TUI.LMem.ScrollDown()
		case "k", "<Up>":
			TUI.LMem.ScrollUp()
		case "g":
			if previousKey == "g" {
				TUI.LMem.ScrollTop()
			}
		case "G", "<End>":
			TUI.LMem.ScrollBottom()
		case "s":
			chip8.EmulateCycle()
			if chip8.DrawFlag {
				TUI.UpdateScreen(&chip8)
				chip8.DrawFlag = false
			}
			TUI.SetEmuInfo(&chip8)
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
	fmt.Printf("| j |Mem map down|\n")
	fmt.Printf("| k | Mem map up |\n")
	fmt.Printf("| gg|Mem map top |\n")
	fmt.Printf("| G |Mem map top |\n")
}

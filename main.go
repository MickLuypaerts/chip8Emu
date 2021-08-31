package main

import (
	"chip8/chip8"
	"fmt"
	"log"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

	// Elements
	lGPR := widgets.NewList()
	lGPR.Title = "registers"
	lGPR.Rows = chip8.GetGPRValues()
	lGPR.TextStyle = ui.NewStyle(ui.ColorYellow)
	lGPR.WrapText = false
	lGPR.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)

	lMem := widgets.NewList()
	lMem.Title = "Memory"
	lMem.TextStyle = ui.NewStyle(ui.ColorYellow)
	lMem.WrapText = false
	lMem.Rows = chip8.GetMemoryValues()

	lProgStats := widgets.NewList()
	lProgStats.Title = "INFO"
	lProgStats.WrapText = true
	lProgStats.Rows = chip8.GetProgStats()

	c := ui.NewCanvas()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(3.0/4, c),          // TODO: program counter, ...
			ui.NewCol(1.0/4, lProgStats), // TODO: Screen
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/4, lGPR),
			ui.NewCol(3.0/4, lMem),
		),
	)
	ui.Render(grid)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			lMem.ScrollDown()
		case "k", "<Up>":
			lMem.ScrollUp()
		case "<C-d>":
			lMem.ScrollHalfPageDown()
		case "<C-u>":
			lMem.ScrollHalfPageUp()
		case "<C-f>":
			lMem.ScrollPageDown()
		case "<C-b>":
			lMem.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				lMem.ScrollTop()
			}
		case "<Home>":
			lMem.ScrollTop()
		case "G", "<End>":
			lMem.ScrollBottom()
		case "r":

		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(grid)
	}
}

func usage() {
	fmt.Printf("Usage: %s [FILE]\n", os.Args[0])
}

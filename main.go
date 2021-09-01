package main

import (
	"chip8/chip8"
	"fmt"
	"image"
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
	lGPR.Title = "Registers"
	lGPR.Rows = chip8.GetGPRValues()
	lGPR.TextStyle = ui.NewStyle(ui.ColorYellow)
	lGPR.WrapText = false
	lGPR.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)

	lStack := widgets.NewList()
	lStack.Title = "Stack"
	lStack.Rows = chip8.GetStackValues()
	lStack.TextStyle = ui.NewStyle(ui.ColorYellow)
	lStack.WrapText = false
	lStack.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)

	lMem := widgets.NewList()
	lMem.Title = "Memory"
	lMem.TextStyle = ui.NewStyle(ui.ColorYellow)
	lMem.WrapText = false
	lMem.Rows = chip8.GetMemoryValues()

	lProgStats := widgets.NewList()
	lProgStats.Title = fmt.Sprintf("INFO %s", os.Args[1])
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
			ui.NewCol(0.5/4, lGPR),
			ui.NewCol(0.5/4, lStack),
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
		case "s":
			chip8.EmulateCycle()
			updateCanvas(c, chip8)
			lProgStats.Rows = chip8.GetProgStats()
			lGPR.Rows = chip8.GetGPRValues()
			lMem.Rows = chip8.GetMemoryValues()
			lStack.Rows = chip8.GetStackValues()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(grid)
	}

	/*
		for {
			var wait string

			fmt.Scanln(&wait)
			chip8.EmulateCycle()

		}*/
}

func usage() {
	fmt.Printf("Usage: %s [FILE]\n", os.Args[0])
}

func updateCanvas(c *ui.Canvas, c8 chip8.Chip8) {
	xOffSet := 2
	yOffSet := 4
	screenBuffer, width, height := c8.GetScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if screenBuffer[x+(y*width)] != 0 {
				c.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorRed)
			} else {
				c.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorWhite)
			}
		}
	}
}

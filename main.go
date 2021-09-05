package main

import (
	"chip8/chip8"
	"chip8/emulator"
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
	uiEvents := ui.PollEvents()
	quit := make(chan struct{})
	controls, err := emulator.CreateKeyFuncMap(chip8.ControlsMap(), TUI.GetControlsMap(), "q", quit)
	if err != nil {
		log.Fatalf("failed to create controls: %v", err)
	}
	go func() {
		<-quit
		emulator.ClearTerminal()
		os.Exit(0)
	}()
	for {
		e := <-uiEvents
		emulator.ExecuteKeyFunction(controls, e.ID)
	}

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

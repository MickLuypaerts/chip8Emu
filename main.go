package main

import (
	"chip8/chip8"
	"chip8/emulator"
	"chip8/view"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

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

	controls, err := emulator.CreateKeyFuncMap(chip8.ControlsMap(), TUI.GetControlsMap())
	if err != nil {
		log.Fatalf("failed to create controls: %v", err)
	}
	quit := make(chan struct{})
	go func() {
		<-quit
		clearTerminal()
		os.Exit(0)
	}()
	controls["q"] = func() { close(quit) }

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

func clearTerminal() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		log.Println("operating system not supported for clearing terminal")
	}
}

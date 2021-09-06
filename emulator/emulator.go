package emulator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Chip interface {
	Init(file string, tuiSetter TUISetter) error
	ControlsMap() map[string]Control

	ChipGetter
}

type ChipGetter interface {
	GetStackValues() []string
	GetScreenSize() (int, int)
	GetGPRValues() []string
	OpcodeInfo() OpcodeInfo
	GetMemoryValues() []byte

	DrawSignal() <-chan []byte
	KeySignal() <-chan []byte
}

type TUI interface {
	Init(drawSignal <-chan []byte, keySignal <-chan []byte, c Chip)
	Close()
	Render()
	Setup() error
	KeyEvent() <-chan string
	ControlsMap() map[string]Control

	TUISetter
}

type TUISetter interface {
	SetEmuInfo(c ChipGetter)
}

var (
	quitSignal = make(chan struct{})
)

type Emulator struct {
	chip     Chip
	tui      TUI
	controls map[string]func()
}

func (emu *Emulator) Run() {
	defer emu.tui.Close()
	emu.tui.Render()

	go func() {
		<-quitSignal
		clearTerminal()
		os.Exit(0)
	}()
	keyEvents := emu.tui.KeyEvent()
	for {
		key := <-keyEvents
		executeKeyFunction(emu.controls, key)
	}
}

func usage(name string, c map[string]Control, t map[string]Control) {
	pUsage, pKey := usagePadding(c, t)
	fmt.Printf("Usage: %s [FILE]\n\n", "chip8")
	fmt.Printf("Emulator Controls:\n")
	fmt.Printf("|" + strings.Repeat("-", pKey+1) + "|" + strings.Repeat("-", pUsage+1) + "|\n")
	fmt.Printf("| %-*s| %-*s|\n", pKey, "key", pUsage, "function")
	fmt.Printf("|" + strings.Repeat("-", pKey+1) + "|" + strings.Repeat("-", pUsage+1) + "|\n")
	for key := range c {
		fmt.Printf("| %-*s| %-*s|\n", pKey, key, pUsage, c[key].usage)
	}
	fmt.Printf("|" + strings.Repeat("-", pKey+1) + "|" + strings.Repeat("-", pUsage+1) + "|\n")
	for key := range t {
		fmt.Printf("| %-*s| %-*s|\n", pKey, key, pUsage, t[key].usage)
	}
}
func usagePadding(c map[string]Control, t map[string]Control) (int, int) {
	maxLenUsage := 0
	maxLenKey := 0
	for key := range c {
		if maxLenUsage < len(c[key].usage) {
			maxLenUsage = len(c[key].usage)
		}
		if maxLenKey < len(key) {
			maxLenKey = len(key)
		}
	}
	for key := range t {
		if maxLenUsage < len(t[key].usage) {
			maxLenUsage = len(t[key].usage)
		}
		if maxLenKey < len(key) {
			maxLenKey = len(key)
		}
	}
	return maxLenUsage + 1, maxLenKey + 1
}

func CreateEmulator(args []string, quitKey string, c Chip, t TUI) (*Emulator, error) {
	e := new(Emulator)
	if len(args) < 2 {
		usage(args[0], c.ControlsMap(), t.ControlsMap())
		os.Exit(0)
	}

	e.chip = c
	e.tui = t
	if err := t.Setup(); err != nil {
		return nil, err
	}

	err := e.chip.Init(args[1], t)
	if err != nil {
		return nil, err
	}

	e.tui.Init(c.DrawSignal(), c.KeySignal(), e.chip)

	e.controls, err = createKeyFuncMap(c.ControlsMap(), t.ControlsMap(), quitKey)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func createKeyFuncMap(chip map[string]Control, tui map[string]Control, quitKey string) (map[string]func(), error) {
	c := make(map[string]func())
	for k, v := range chip {
		if _, ok := c[k]; !ok {
			c[k] = v.f
		} else {
			return nil, DoubleKeyAssigmentError{Key: k}
		}
	}
	for k, v := range tui {
		if _, ok := c[k]; !ok {
			c[k] = v.f
		} else {
			return nil, DoubleKeyAssigmentError{Key: k}
		}
	}
	if _, ok := c[quitKey]; !ok {
		c[quitKey] = func() { close(quitSignal) }
	}
	return c, nil
}

func executeKeyFunction(m map[string]func(), k string) {
	if f, ok := m[k]; ok {
		f()
	}
}

func clearTerminal() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		log.Println("operating system not supported for clearing terminal")
	}
}

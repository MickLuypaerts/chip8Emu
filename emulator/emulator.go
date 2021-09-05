package emulator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	ui "github.com/gizak/termui/v3"
)

type Chip interface {
	ChipIniter

	GetKeyValues() []string
	GetStackValues() []string

	GetScreen() []byte
	GetScreenSize() (int, int)
	GetGPRValues() []string
	OpcodeInfo() OpcodeInfo
	GetMemoryValues() []byte
}

type ChipIniter interface {
	Init(file string, screenFunc func(Chip), infoFunc func(Chip), keyFunc func(Chip)) error
	ControlsMap() map[string]Control
}

type TUI interface {
	TUIIniter

	Grid() *ui.Grid
	UpdateScreen(c Chip)
	SetEmuInfo(c Chip)
	SetKeyInfo(c Chip)
}

type TUIIniter interface {
	Init(c Chip)
	ControlsMap() map[string]Control
}

type Emulator struct {
	chip     Chip
	tui      TUI
	controls map[string]func()

	quitKey string
	quit    chan struct{}
}

func (emu *Emulator) Run() {
	defer ui.Close()
	ui.Render(emu.tui.Grid())

	go func() {
		<-emu.quit
		clearTerminal()
		os.Exit(0)
	}()
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		executeKeyFunction(emu.controls, e.ID)
	}
}

func usage(name string, c map[string]Control, t map[string]Control) {
	// TODO: usage()
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

func CreateEmulator(args []string, quitKey string, c Chip, t TUI) (*Emulator, error) {
	e := new(Emulator)
	e.quitKey = quitKey
	e.quit = make(chan struct{})
	if len(args) < 2 {
		usage(args[0], c.ControlsMap(), t.ControlsMap())
		os.Exit(0)
	}
	if err := ui.Init(); err != nil {
		return nil, err
	}
	e.chip = c
	err := e.chip.Init(args[1], t.UpdateScreen, t.SetEmuInfo, t.SetKeyInfo)
	if err != nil {
		return nil, err
	}
	e.tui = t
	e.tui.Init(e.chip)

	e.controls, err = createKeyFuncMap(c.ControlsMap(), t.ControlsMap(), quitKey, e.quit)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func createKeyFuncMap(chip map[string]Control, tui map[string]Control, quitKey string, quitChan chan struct{}) (map[string]func(), error) {
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
		c[quitKey] = func() { close(quitChan) }
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

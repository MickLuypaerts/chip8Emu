package emulator

import "os"

type Chip interface {
	ChipIniter
}

type ChipIniter interface {
	Init(file string) error
	ControlsMap() map[string]Control
}

type TUI interface {
	TUIIniter
}

type TUIIniter interface {
	Init() error
	ControlsMap() map[string]Control
}

type Emulator struct {
	chip     Chip
	tui      TUI
	controls map[string]func()
}

func CreateEmulator(args []string, c Chip, t TUI) (*Emulator, error) {
	e := new(Emulator)
	if len(args) < 2 {
		usage(args[0], c.ControlsMap(), t.ControlsMap())
		os.Exit(0)
	}
	e.chip = c
	err := e.chip.Init(args[1])
	if err != nil {
		return nil, err
	}
	e.tui = t
	err = e.tui.Init()
	if err != nil {
		return nil, err
	}
	e.controls, err = createKeyFuncMap(c.ControlsMap(), t.ControlsMap())
	if err != nil {
		return nil, err
	}
	return e, nil
}

func usage(name string, c map[string]Control, t map[string]Control) {
	// TODO: usage()
}

func createKeyFuncMap(chip map[string]Control, tui map[string]Control) (map[string]func(), error) {
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
	return c, nil
}

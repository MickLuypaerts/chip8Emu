package view

import "chip8/emulator"

func (t TUI) GetControlsMap() map[string]emulator.Control {
	m := make(map[string]emulator.Control)
	m["j"] = emulator.NewControl(func() { scrollDown(t.lMem) }, "Mem map down")
	m["<Down>"] = emulator.NewControl(func() { scrollDown(t.lMem) }, "Mem map down")
	m["k"] = emulator.NewControl(func() { scrollUp(t.lMem) }, "Mem map up")
	m["<Up>"] = emulator.NewControl(func() { scrollUp(t.lMem) }, "Mem map up")
	m["g"] = emulator.NewControl(func() { scrollTop(t.lMem) }, "Mem map top")
	m["G"] = emulator.NewControl(func() { scrollBottom(t.lMem) }, "Mem map bottom")
	return m
}

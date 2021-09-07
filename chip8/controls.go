package chip8

import "chip8/emulator"

func (c Chip8) ControlsMap() map[string]emulator.Control {
	m := make(map[string]emulator.Control)
	m["1"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x1) }, "")
	m["2"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x2) }, "")
	m["3"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x3) }, "")
	m["4"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x4) }, "")
	m["5"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x5) }, "")
	m["6"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x6) }, "")
	m["7"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x7) }, "")
	m["8"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x8) }, "")
	m["9"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0x9) }, "")
	m["A"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0xA) }, "")
	m["B"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0xB) }, "")
	m["C"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0xC) }, "")
	m["D"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0xD) }, "")
	m["E"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0xE) }, "")
	m["F"] = emulator.NewControl(func() { sendKeyboardInterrupt(keyboardInterrupt, 0xF) }, "")

	m["r"] = emulator.NewControl(c.run, "run rom")
	m["R"] = emulator.NewControl(c.stop, "stop rom")
	m["s"] = emulator.NewControl(c.emulateCycle, "run 1 cycle")
	return m
}

func sendKeyboardInterrupt(c chan byte, key byte) {
	if running {
		c <- key
	}
}

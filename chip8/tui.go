package chip8

import (
	"chip8/emulator"
	"fmt"
)

func (c *Chip8) GetGPRValues() []string {
	var gpr []string
	for i := range c.v {
		if c.vChanged[i] {
			gpr = append(gpr, fmt.Sprintf("V%X: | %X |", i, c.v[i]))
			c.vChanged[i] = false
		} else {
			gpr = append(gpr, fmt.Sprintf("V%X:   %X", i, c.v[i]))
		}
	}
	gpr = append(gpr, fmt.Sprintf("I :   %03X", c.i))
	return gpr
}

func (c Chip8) GetStackValues() []string {
	var stack []string

	for i := range c.stack {
		stack = append(stack, fmt.Sprintf("%X: 0x%04X", i, c.stack[i]))
	}
	stack = append(stack, fmt.Sprintf("SP: %X", c.sp))
	return stack
}

func (c Chip8) GetMemoryValues() []byte {
	return c.memory[:]
}

func (c Chip8) GetScreenSize() (int, int) {
	return screenWidth, screenHeigth
}

func (c Chip8) OpcodeInfo() emulator.OpcodeInfo {
	return c.info
}

package chip8

import "fmt"

func (c *Chip8) GetKeyValues() []string {
	var keys []string
	for i := range c.key {
		keys = append(keys, fmt.Sprintf("K%X   %d", i, c.key[i]))
	}
	return keys
}

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

const (
	lMemRowLength = 16
)

func (c Chip8) GetMemoryValues() []string {
	var mem []string
	var row string

	rowCount := 0
	columnCount := 0x0

	for i := 0; i < len(c.memory); i++ {
		if rowCount == 0 {
			row = fmt.Sprintf("0x%04X ", columnCount<<4)
			columnCount++
		}
		row += fmt.Sprintf("%02X ", c.memory[i])
		rowCount++
		if rowCount == lMemRowLength {
			mem = append(mem, row)
			rowCount = 0
		}
	}
	return mem
}

func (c Chip8) GetMemoryRow() uint16 {
	return uint16(c.pc / lMemRowLength)
}

func (c Chip8) GetScreen() ([]byte, int, int) {
	screen := c.screenBuf[:]
	return screen, screenWidth, screenHeigth
}

func (c Chip8) GetProgStats() []string {
	progStats := []string{
		fmt.Sprintf("OPCODE: 0x%04X", c.opcode),
		fmt.Sprintf("Name:     %s", c.info.opcodeName),
		fmt.Sprintf("Type: %s", c.info.opcodeType),
		fmt.Sprintf("Desc: %s", c.info.opcodeDesc),
		fmt.Sprintf("PC: %d", c.pc),
		fmt.Sprintf("Index: %04X", c.i),
		fmt.Sprintf("Sound: %v", c.info.playingSound),
	}
	return progStats
}

package chip8

import (
	"fmt"
	"log"
)

func (c *Chip8) decode() {
	switch c.opcode & 0xF000 {
	case 0x0000:
		c.decode0x000()
	case 0x1000:
		c.setOpcodeInfo("1NNN", "Flow", "Jumps to address NNN.")
		c.pc = c.opcode & 0x0FFF
	case 0x2000:
		c.setOpcodeInfo("2NNN", "Flow", "Calls subroutine at NNN.")

		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.opcode & 0x0FFF
	case 0x3000:
		c.setOpcodeInfo("3XNN", "Cond", "Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block);")
		index := c.opcode & 0x0F00 >> 8
		if c.v[index] == byte(c.opcode&0x00FF) {
			c.pc += 2
		}
		c.pc += 2

	case 0x4000:
		c.setOpcodeInfo("4XNN", "Cond", "Skips the next instruction if VX does not equal NN. (Usually the next instruction is a jump to skip a code block);")
		if c.v[c.opcode&0x0F00>>8] != byte(c.opcode&0x00FF) {
			c.pc += 2
		}
		c.pc += 2
	case 0x5000:
		c.setOpcodeInfo("5XY0", "Cond", "Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block);")
		if c.v[c.opcode&0x0F00>>8] == c.v[c.opcode&0x00F0>>4] {
			c.pc += 2
		}
		c.pc += 2

	case 0x6000: // 6xkk
		c.setOpcodeInfo("6XNN", "Const", "Sets VX to NN.")
		index := c.opcode & 0x0F00 >> 8
		c.v[index] = byte(c.opcode & 0x00FF)
		c.vChanged[index] = true
		c.pc += 2
	case 0x7000:
		c.setOpcodeInfo("7XNN", "Const", "Adds NN to VX. (Carry flag is not changed);")
		index := (c.opcode & 0x0F00) >> 8
		c.v[index] += byte(c.opcode & 0x00FF)
		c.vChanged[index] = true
		c.pc += 2
		// 7xkk
	case 0x8000:
		// 8xy0
		// 8xy1
		// 8xy2
		// 8xy3
		// 8xy4
		// 8xy5
		// 8xy6
		// 8xy7
		// 8xyE
	case 0x9000:
		c.setOpcodeInfo("9XY0", "Cond", "Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block);")
		if c.v[c.opcode&0x0F00>>8] != c.v[c.opcode&0x00F0>>4] {
			c.pc += 2
		}
		c.pc += 2
	case 0xA000:
		c.setOpcodeInfo("ANNN", "MEM", "Sets I to the address NNN.")
		c.i = c.opcode & 0x0FFF
		c.pc += 2
	case 0xB000:
		// Bnnn
	case 0xC000:
		// Cxkk
	case 0xD000:
		c.setOpcodeInfo("DXYN", "Disp", "Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.")
		x := uint16(c.v[(c.opcode&0x0F00)>>8])
		y := uint16(c.v[(c.opcode&0x00F0)>>4])
		h := uint16(c.opcode & 0x000F)
		c.draw(x, y, h)
		c.pc += 2
		// Dxyn
	case 0xE000:
		// Ex9E
		// ExA1
	case 0xF000:
		// Fx07
		// Fx0A
		// Fx15
		// Fx18
		// Fx1E
		// Fx29
		// Fx33
		// Fx55
		// Fx65
	default:
		log.Printf("[ERROR]: Unknown opcode: ox%X\n", c.opcode)
	}
}

func (c *Chip8) decode0x000() {
	// 0nnn
	// 00E0
	// 00EE
	switch c.opcode & 0x00FF {
	case 0x00E0:
		c.setOpcodeInfo("00E0", "Display", "Clears the screen.")
		c.clearScreen()
		c.pc += 2
	case 0x00EE:
		c.setOpcodeInfo("00EE", "Flow", "Returns from a subroutine.")
		c.sp--
		c.pc = c.stack[c.sp]
	default:
		c.setOpcodeInfo("0NNN", "Call", "Calls machine code routine. Not implemented")
		c.pc += 2
	}
}

func (c *Chip8) draw(x, y, h uint16) {
	var pixel uint16
	for yLine := 0; yLine < int(h); yLine++ {
		pixel = uint16(c.memory[c.i+uint16(yLine)]) // Fetch the pixel value from the memory starting at location I
		for xLine := 0; xLine < 8; xLine++ {
			if (pixel & (0x80 >> xLine)) != 0 {
				if c.screenBuf[(x+uint16(xLine)+((y+uint16(yLine))*screenWidth))] == 1 { // Check if the pixel on the display is set to 1. If it is set,
					c.v[0xF] = 1 // we need to register the collision by setting the VF register
				}
				c.screenBuf[x+uint16(xLine)+((y+uint16(yLine))*screenWidth)] ^= 1
			}

		}
	}
	c.DrawFlag = true
	c.vChanged[0xF] = true
	//c.printScreenToConsole()
}

func (c Chip8) printScreenToConsole() {
	for y := 0; y < screenHeigth; y++ {
		for x := 0; x < screenWidth; x++ {
			fmt.Print(c.screenBuf[x+(y*screenWidth)])
		}
		fmt.Println()
	}
	fmt.Println()
}

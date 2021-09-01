package chip8

import (
	"fmt"
	"log"
)

func (c *Chip8) decode() {
	c.pc += 2
	switch c.opcode & 0xF000 {
	case 0x0000:
		c.decode0x0000()
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
		// index := c.opcode & 0x0F00 >> 8
		if c.v[c.getXFromOpcode()] == byte(c.opcode&0x00FF) {
			c.pc += 2
		}

	case 0x4000:
		c.setOpcodeInfo("4XNN", "Cond", "Skips the next instruction if VX does not equal NN. (Usually the next instruction is a jump to skip a code block);")
		if c.v[c.getXFromOpcode()] != byte(c.opcode&0x00FF) {
			c.pc += 2
		}
	case 0x5000:
		c.setOpcodeInfo("5XY0", "Cond", "Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block);")
		if c.v[c.getXFromOpcode()] == c.v[c.getYFromOpcode()] {
			c.pc += 2
		}

	case 0x6000: // 6xkk
		c.setOpcodeInfo("6XNN", "Const", "Sets VX to NN.")
		index := c.getXFromOpcode()
		c.v[index] = byte(c.opcode & 0x00FF)
		c.vChanged[index] = true
	case 0x7000:
		c.setOpcodeInfo("7XNN", "Const", "Adds NN to VX. (Carry flag is not changed);")
		// index := (c.opcode & 0x0F00) >> 8
		index := c.getXFromOpcode()
		c.v[index] += byte(c.opcode & 0x00FF)
		c.vChanged[index] = true
		// 7xkk
	case 0x8000:
		c.decode0x8000()
	case 0x9000:
		c.setOpcodeInfo("9XY0", "Cond", "Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block);")
		if c.v[c.getXFromOpcode()] != c.v[c.getYFromOpcode()] {
			c.pc += 2
		}
	case 0xA000:
		c.setOpcodeInfo("ANNN", "MEM", "Sets I to the address NNN.")
		c.i = c.opcode & 0x0FFF
	case 0xB000:
		// Bnnn
		c.pc -= 2
	case 0xC000:
		// Cxkk
		c.pc -= 2
	case 0xD000:
		c.setOpcodeInfo("DXYN", "Disp", "Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.")
		x := uint16(c.v[c.getXFromOpcode()])
		y := uint16(c.v[c.getYFromOpcode()])
		h := uint16(c.opcode & 0x000F)
		c.draw(x, y, h)
		// Dxyn
	case 0xE000:
		c.pc -= 2
		// Ex9E
		// ExA1
	case 0xF000:
		c.decode0xF000()
	default:
		log.Printf("[ERROR]: Unknown opcode: ox%X\n", c.opcode)
	}
}
func (c *Chip8) decode0xF000() {
	switch c.opcode & 0x00FF {
	case 0x0007: // Fx07
		c.setOpcodeInfo("FX07", "Timer", "Sets VX to the value of the delay timer.")
		c.v[c.getXFromOpcode()] = c.delayTimer
	case 0x000A: // Fx0A
		c.pc -= 2
	case 0x0015: // Fx15
		c.setOpcodeInfo("FX15", "Timer", "Sets the delay timer to VX.")
		c.delayTimer = c.v[c.getXFromOpcode()]
	case 0x0018: // Fx18
		c.setOpcodeInfo("FX18", "Sound", "Sets the sound timer to VX.")
		c.soundTimer = c.v[c.getXFromOpcode()]
	case 0x001E: // Fx1E
		c.pc -= 2
	case 0x0029: // Fx29
		c.pc -= 2
	case 0x0033: // Fx33
		c.pc -= 2
	case 0x0055: // TODO: FX55 should we increment I here?
		c.setOpcodeInfo("FX55", "MEM", "Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.")
		for i := 0; i <= int(c.getXFromOpcode()); i++ {
			c.memory[c.i] = c.v[i]
			c.i++
		}
		c.i++
	case 0x0065: // Fx65
		c.pc -= 2
	}
}

func (c *Chip8) decode0x8000() {
	switch c.opcode & 0x000F {
	case 0x0000:
		c.setOpcodeInfo("8XY0", "Assig", "Sets VX to the value of VY.")
		index := c.getXFromOpcode()
		c.v[index] = c.v[c.opcode&0x00F0>>4]
		c.vChanged[index] = true
	case 0x0001:
		c.setOpcodeInfo("8XY1", "BitOp", "Sets VX to VX or VY. (Bitwise OR operation);")
		index := c.getXFromOpcode()
		c.v[index] |= c.v[c.opcode&0x00F0>>4]
		c.vChanged[index] = true
	case 0x0002:
		c.setOpcodeInfo("8XY2", "BitOp", "Sets VX to VX and VY. (Bitwise AND operation);")
		index := c.getXFromOpcode()
		c.v[index] &= c.v[c.opcode&0x00F0>>4]
		c.vChanged[index] = true
	case 0x0003:
		c.setOpcodeInfo("8XY3", "BitOp", "Sets VX to VX xor VY. (Bitwise XOR operation);")
		index := c.getXFromOpcode()
		c.v[index] ^= c.v[c.opcode&0x00F0>>4]
		c.vChanged[index] = true
	case 0x0004:
		c.setOpcodeInfo("8XY4", "Math", "Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.")
		indexX := c.getXFromOpcode()
		indexY := c.getYFromOpcode()
		if c.v[indexX] > (0xFF - c.v[indexY]) {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.vChanged[0xF] = true
		c.v[indexX] += c.v[indexY]
		c.vChanged[indexX] = true
	case 0x0005: // TODO: double check 8XY5
		c.setOpcodeInfo("8XY5", "Math", "VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.")
		indexX := c.getXFromOpcode()
		indexY := c.getYFromOpcode()
		c.subtract(indexX, indexX, indexY)
	case 0x0006:
		c.setOpcodeInfo("8XY6", "BitOp", "Stores the least significant bit of VX in VF and then shifts VX to the right by 1.")
		index := c.getXFromOpcode()
		if c.v[index]&0x01 == 1 {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.vChanged[0xF] = true
		c.v[index] <<= 1
		c.vChanged[index] = true
	case 0x0007:
		// 8xy7
		c.setOpcodeInfo("8XY7", "Math", "Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not.")
		indexX := c.getXFromOpcode()
		indexY := c.getYFromOpcode()
		c.subtract(indexX, indexY, indexX)
	case 0x000E:
		c.setOpcodeInfo("8XYE", "BitOp", "Stores the most significant bit of VX in VF and then shifts VX to the left by 1.")
		index := c.getXFromOpcode()
		if c.v[index]&0x8 == 1 {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.vChanged[0xF] = true
		c.v[index] >>= 1
		c.vChanged[index] = true
	default:
		log.Printf("[ERROR]: Unknown opcode: ox%X\n", c.opcode)

	}
}

func (c *Chip8) decode0x0000() {
	switch c.opcode & 0x00FF {
	case 0x00E0:
		c.setOpcodeInfo("00E0", "Display", "Clears the screen.")
		c.clearScreen()
	case 0x00EE:
		c.setOpcodeInfo("00EE", "Flow", "Returns from a subroutine.")
		c.sp--
		c.pc = c.stack[c.sp]
	default:
		c.setOpcodeInfo("0NNN", "Call", "Calls machine code routine. Not implemented")
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
					c.vChanged[0xF] = true
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

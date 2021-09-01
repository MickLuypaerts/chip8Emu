package chip8

import "log"

func (c *Chip8) decode0x000() {
	// 0nnn
	// 00E0
	// 00EE
	switch c.opcode & 0x00FF {
	case 0x00E0: // 00E0
		c.setOpcodeInfo("00E0", "Display", "Clears the screen.")
		c.clearScreen()
		c.pc += 2
	case 0x00EE: // 00EE
		c.setOpcodeInfo("00EE", "Flow", "Returns from a subroutine.")
		c.pc = c.stack[c.sp]
		c.sp--
	default:
		c.setOpcodeInfo("0NNN", "Call", "Calls machine code routine. Not implemented")
		c.pc += 2 // 0nnn
	}
}

func (c *Chip8) decode() {
	switch c.opcode & 0xF000 {
	case 0x0000:
		c.decode0x000()
	case 0x1000:
		// 1nnn
	case 0x2000:
		// 2nnn
	case 0x3000:
		// 3xkk
	case 0x4000:
		// 4xkk
	case 0x5000:
		// 5xy0
	case 0x6000:
		// 6xkk
	case 0x7000:
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
		// 9xy0
	case 0xA000:
		// Annn
	case 0xB000:
		// Bnnn
	case 0xC000:
		// Cxkk
	case 0xD000:
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

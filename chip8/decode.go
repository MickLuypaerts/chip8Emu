package chip8

import (
	"log"
	"math/rand"
	"time"
)

type opcodeParts struct {
	x   byte
	y   byte
	n   byte
	nn  byte
	nnn uint16
}

// TODO: should decode just return the opcode name
//       and then have a method Execute that uses a map[string]func() to execute the opcode

// TODO: remove all the c.vChanged[index] = true from here to somewhere else change of forgetting to add it
//       are too big
func (c *Chip8) decode() {
	c.pc += 2
	o := opcodeParts{x: xFromOpcode(c.opcode), y: yFromOpcode(c.opcode), nnn: nnnFromOpcode(c.opcode), nn: nnFromOpcode(c.opcode), n: nFromOpcode(c.opcode)}
	switch c.opcode & 0xF000 {
	case 0x0000:
		c.decode0x0000(o)
	case 0x1000:
		c.pc = o.nnn
		c.setEmulatorInfo("1NNN", "Flow", "Jumps to address NNN.")
	case 0x2000:
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = o.nnn
		c.setEmulatorInfo("2NNN", "Flow", "Calls subroutine at NNN.")
	case 0x3000:
		if c.v[o.x] == byte(o.nn) {
			c.pc += 2
		}
		c.setEmulatorInfo("3XNN", "Cond", "Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block);")
	case 0x4000:
		if c.v[o.x] != o.nn {
			c.pc += 2
		}
		c.setEmulatorInfo("4XNN", "Cond", "Skips the next instruction if VX does not equal NN. (Usually the next instruction is a jump to skip a code block);")
	case 0x5000:
		if c.v[o.x] == c.v[o.y] {
			c.pc += 2
		}
		c.setEmulatorInfo("5XY0", "Cond", "Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block);")
	case 0x6000:
		c.v[o.x] = o.nn
		c.vChanged[o.x] = true
		c.setEmulatorInfo("6XNN", "Const", "Sets VX to NN.")
	case 0x7000:
		c.v[o.x] += o.nn
		c.vChanged[o.x] = true
		c.setEmulatorInfo("7XNN", "Const", "Adds NN to VX. (Carry flag is not changed);")
	case 0x8000:
		c.decode0x8000(o)
	case 0x9000:
		if c.v[o.x] != c.v[o.y] {
			c.pc += 2
		}
		c.setEmulatorInfo("9XY0", "Cond", "Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block);")
	case 0xA000:
		c.i = o.nnn
		c.setEmulatorInfo("ANNN", "MEM", "Sets I to the address NNN.")
	case 0xB000:
		c.pc = o.nnn + uint16(c.v[0x0])
		c.setEmulatorInfo("BNNN", "Flow", "Jumps to the address NNN plus V0.")
	case 0xC000:
		randSource := rand.NewSource(time.Now().UnixNano())
		r := rand.New(randSource)
		c.v[o.x] = byte(r.Intn(256)) & o.nn
		c.vChanged[o.x] = true
		c.setEmulatorInfo("CXNN", "Rand", "Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.")
	case 0xD000:
		x := uint16(c.v[o.x])
		y := uint16(c.v[o.y])
		h := uint16(o.n)
		c.draw(x, y, h)
		c.setEmulatorInfo("DXYN", "Disp", "Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.")
	case 0xE000:
		c.decode0xE000(o)
	case 0xF000:
		c.decode0xF000(o)
	default:
		log.Printf("[ERROR]: Unknown opcode: ox%X\n", c.opcode)
	}
}
func (c *Chip8) decode0xE000(o opcodeParts) {
	switch c.opcode & 0x0FF {
	case 0x009E:
		if c.key[c.v[o.x]] == 1 {
			c.pc += 2
		}
		c.setEmulatorInfo("EX9E", "KeyOp", "Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block);")
	case 0x00A1:
		if c.key[c.v[o.x]] != 1 {
			c.pc += 2
		}
		c.setEmulatorInfo("EXA1", "KeyOp", "Skips the next instruction if the key stored in VX is not pressed. (Usually the next instruction is a jump to skip a code block);")
	}
}

func (c *Chip8) decode0xF000(o opcodeParts) {
	switch c.opcode & 0x00FF {
	case 0x0007:
		c.v[o.x] = c.delayTimer
		c.setEmulatorInfo("FX07", "Timer", "Sets VX to the value of the delay timer.")
	case 0x000A:
		// FX0A  KeyOp  A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event);
		// just wait for keyboardInterrupt here?
		c.pc -= 2
	case 0x0015:
		c.delayTimer = c.v[o.x]
		c.setEmulatorInfo("FX15", "Timer", "Sets the delay timer to VX.")
	case 0x0018:
		c.soundTimer = c.v[o.x]
		c.setEmulatorInfo("FX18", "Sound", "Sets the sound timer to VX.")
	case 0x001E:
		c.i += uint16(c.v[o.x])
		c.setEmulatorInfo("FX1E", "MEM", "Adds VX to I. VF is not affected.")
	case 0x0029:
		var loc uint16
		for i := byte(0x0); i < 0x10; i++ {
			if c.v[o.x] == i {
				c.i = loc
			}
			loc += 5
		}
		c.setEmulatorInfo("FX29", "MEM", "Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.")
	case 0x0033:
		/*
			Store BCD representation of Vx in I
			take decimal number of Vx and place
			X00 = I
			0X0 = I+1
			00X = I+2
		*/
		c.memory[c.i] = c.v[o.x] / 100
		c.memory[c.i+1] = (c.v[o.x] / 10) % 10
		c.memory[c.i+2] = (c.v[o.x] % 100) % 10
		c.setEmulatorInfo("FX33", "BCD", "Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.);")
	case 0x0055: // TODO: FX55 should we increment I here?
		for i := byte(0x0); i <= o.x; i++ {
			c.memory[c.i] = c.v[i]
			c.i++
		}
		c.i++
		c.setEmulatorInfo("FX55", "MEM", "Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.")
	case 0x0065: // TODO: FX65 should we increment I here?
		for i := byte(0x0); i <= o.x; i++ {
			c.v[i] = c.memory[c.i]
			c.vChanged[i] = true
			c.i++
		}
		c.i++
		c.setEmulatorInfo("FX65", "MEM", "Fills V0 to VX (including VX) with values from memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.")
	}
}

func (c *Chip8) decode0x8000(o opcodeParts) {
	switch c.opcode & 0x000F {
	case 0x0000:
		c.v[o.x] = c.v[o.y]
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XY0", "Assig", "Sets VX to the value of VY.")
	case 0x0001:
		c.v[o.x] |= c.v[o.y]
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XY1", "BitOp", "Sets VX to VX or VY. (Bitwise OR operation);")
	case 0x0002:
		c.v[o.x] &= c.v[o.y]
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XY2", "BitOp", "Sets VX to VX and VY. (Bitwise AND operation);")
	case 0x0003:
		c.v[o.x] ^= c.v[o.y]
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XY3", "BitOp", "Sets VX to VX xor VY. (Bitwise XOR operation);")
	case 0x0004:
		if c.v[o.x] > (0xFF - c.v[o.y]) {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.vChanged[0xF] = true
		c.v[o.x] += c.v[o.y]
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XY4", "Math", "Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.")
	case 0x0005: // TODO: double check 8XY5
		c.subtract(o.x, o.x, o.y)
		c.setEmulatorInfo("8XY5", "Math", "VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.")
	case 0x0006:
		if c.v[o.x]&0x01 == 1 {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.vChanged[0xF] = true
		c.v[o.x] >>= 1
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XY6", "BitOp", "Stores the least significant bit of VX in VF and then shifts VX to the right by 1.")
	case 0x0007:
		c.subtract(o.x, o.y, o.x)
		c.setEmulatorInfo("8XY7", "Math", "Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not.")
	case 0x000E:
		if c.v[o.x]&0x8 == 1 {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.vChanged[0xF] = true
		c.v[o.x] <<= 1
		c.vChanged[o.x] = true
		c.setEmulatorInfo("8XYE", "BitOp", "Stores the most significant bit of VX in VF and then shifts VX to the left by 1.")
	default:
		log.Printf("[ERROR]: Unknown opcode: ox%X\n", c.opcode)

	}
}

func (c *Chip8) decode0x0000(o opcodeParts) {
	switch c.opcode & 0x00FF {
	case 0x00E0:
		c.clearScreen()
		c.setEmulatorInfo("00E0", "Display", "Clears the screen.")
	case 0x00EE:
		c.sp--
		c.pc = c.stack[c.sp]
		c.setEmulatorInfo("00EE", "Flow", "Returns from a subroutine.")
	default:
		c.setEmulatorInfo("0NNN", "Call", "Calls machine code routine. Not implemented")
	}
}

func (c *Chip8) draw(x, y, h uint16) {
	var pixel uint16
	for yLine := 0; yLine < int(h); yLine++ {
		pixel = uint16(c.memory[c.i+uint16(yLine)]) // Fetch the pixel value from the memory starting at location I
		for xLine := 0; xLine < 8; xLine++ {
			if (pixel & (0x80 >> xLine)) != 0 {
				index := (x + uint16(xLine) + ((y + uint16(yLine)) * screenWidth))
				if index < uint16(len(c.screenBuf)) {
					if c.screenBuf[index] == 1 { // Check if the pixel on the display is set to 1. If it is set,
						c.v[0xF] = 1 // we need to register the collision by setting the VF register
						c.vChanged[0xF] = true
					}
					c.screenBuf[x+uint16(xLine)+((y+uint16(yLine))*screenWidth)] ^= 1
				}
			}
		}
	}
	c.drawFlag = true
	c.vChanged[0xF] = true
}

func (c *Chip8) clearScreen() {
	for i := range c.screenBuf {
		c.screenBuf[i] = 0
	}
}

func (c *Chip8) subtract(target, x, y byte) {
	if c.v[x] > c.v[y] {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
	c.vChanged[0xF] = true
	c.v[target] = c.v[x] - c.v[y]
	c.vChanged[target] = true
}

func xFromOpcode(opcode uint16) byte {
	return byte((opcode & 0x0F00) >> 8)
}

func yFromOpcode(opcode uint16) byte {
	return byte((opcode & 0x00F0) >> 4)
}

func nnnFromOpcode(opcode uint16) uint16 {
	return opcode & 0x0FFF
}

func nnFromOpcode(opcode uint16) byte {
	return byte(opcode & 0x00FF)
}

func nFromOpcode(opcode uint16) byte {
	return byte(opcode & 0x000F)
}

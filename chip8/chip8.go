package chip8

import (
	"io/ioutil"
)

const (
	memorySize   = 4096
	vRegSize     = 16
	stackSize    = 16
	screenWidth  = 64
	screenHeigth = 32
	keyNumbers   = 16
)

type Chip8 struct {
	opcode     uint16
	i          uint16 // The address register, which is named I, is 12 bits wide and is used with several opcodes that involve memory operations.
	pc         uint16
	stack      [stackSize]uint16
	sp         byte
	memory     [memorySize]byte
	v          [vRegSize]byte // general purpose registers
	vChanged   [vRegSize]bool
	screenBuf  [screenWidth * screenHeigth]byte
	key        [keyNumbers]byte
	delayTimer byte
	soundTimer byte

	info chip8Info
}

type chip8Info struct {
	playingSound bool
}

func (c *Chip8) Init(file string) error {
	c.opcode = 0x200 // programs written for the original system begin at memory location 512 (0x200)
	romData, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	// load fontset
	fontset := [80]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	for i := range fontset {
		c.memory[i] = fontset[i]
	}

	// load program into memory
	for i := range romData {
		c.memory[i+512] = romData[i]
	}
	return nil
}

func (c *Chip8) clearScreen() {
	for i := range c.screenBuf {
		c.screenBuf[i] = 0
	}
}

func (c *Chip8) fetch() {
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
}

func (c *Chip8) updateTimers() {
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.info.playingSound = true
		c.soundTimer--
	} else {
		c.info.playingSound = false
	}
}

func (c *Chip8) EmulateCycle() {
	c.fetch()
	c.decode()
	c.updateTimers()
}

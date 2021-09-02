package chip8

import (
	"chip8/view"
	"io/ioutil"
	"time"
)

const (
	memorySize     = 4096
	vRegSize       = 16
	stackSize      = 16
	screenWidth    = 64
	screenHeigth   = 32
	keyNumbers     = 16
	clockCycleRate = 2 * time.Microsecond
	timeCycleRate  = 16 * time.Microsecond
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
	drawFlag   bool
	key        [keyNumbers]byte
	delayTimer byte
	soundTimer byte
	stop       chan bool

	info       chip8Info
	ScreenFunc func(view.Chip)
	InfoFunc   func(view.Chip)
}

type chip8Info struct {
	opcodeName   string
	opcodeType   string
	opcodeDesc   string
	playingSound bool
}

func (c *Chip8) setOpcodeInfo(n string, t string, d string) {
	c.info.opcodeName = n
	c.info.opcodeType = t
	c.info.opcodeDesc = d
}

func (c *Chip8) Init(file string) error {
	c.pc = 0x200 // programs written for the original system begin at memory location 512 (0x200)
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
	c.stop = make(chan bool)
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
	if c.drawFlag {
		c.ScreenFunc(c)
		c.drawFlag = false
	}
	c.InfoFunc(c)
}

func (c *Chip8) Run() {
	clock := time.NewTicker(clockCycleRate)
	timers := time.NewTicker(timeCycleRate)

	go c.runCycle(c.EmulateCycle, clock)
	go c.runCycle(c.updateTimers, timers)
}

func (c *Chip8) Stop() {
	c.stop <- true
	c.stop <- true
}

func (c *Chip8) runCycle(f func(), cycle *time.Ticker) {
	go func() {
		for {
			select {
			case <-c.stop:
				cycle.Stop()
				return
			case <-cycle.C:
				f()
			}
		}
	}()
}
func (c *Chip8) subtract(target, x, y uint16) {
	if c.v[x] > c.v[y] {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
	c.vChanged[0xF] = true
	c.v[target] = c.v[x] - c.v[y]
	c.vChanged[target] = true
}

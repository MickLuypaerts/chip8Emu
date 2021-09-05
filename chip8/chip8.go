package chip8

import (
	"chip8/view"
	"io/ioutil"
	"time"
)

const (
	memorySize        = 4096
	vRegSize          = 16
	stackSize         = 16
	screenWidth       = 64
	screenHeigth      = 32
	keyNumbers        = 16
	clockCycleRate    = 2 * time.Microsecond
	timeCycleRate     = 16 * time.Microsecond
	keyboardCycleRate = 8 * time.Millisecond
)

type Chip8 struct {
	opcode            uint16
	i                 uint16 // The address register, which is named I, is 12 bits wide and is used with several opcodes that involve memory operations.
	pc                uint16
	stack             [stackSize]uint16
	sp                byte
	memory            [memorySize]byte
	v                 [vRegSize]byte // general purpose registers
	vChanged          [vRegSize]bool
	screenBuf         [screenWidth * screenHeigth]byte
	drawFlag          bool
	key               [keyNumbers]byte
	delayTimer        byte
	soundTimer        byte
	stopSignal        chan bool
	keyBoardInterrupt chan byte

	info       chip8Info
	ScreenFunc func(view.Chip)
	InfoFunc   func(view.Chip)
	KeyFunc    func(view.Chip)
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

func (c *Chip8) Init(file string, screenFunc func(view.Chip), infoFunc func(view.Chip), keyFunc func(view.Chip)) error {
	c.pc = 0x200 // programs written for the original system begin at memory location 512 (0x200)
	c.ScreenFunc = screenFunc
	c.InfoFunc = infoFunc
	c.KeyFunc = keyFunc
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
	c.stopSignal = make(chan bool)
	c.keyBoardInterrupt = make(chan byte)
	return nil
}

func (c *Chip8) ClearKeys() {
	clearedKey := false
	for i := range c.key {
		if c.key[i] != 0 {
			c.key[i] = 0
			clearedKey = true
		}
	}
	if clearedKey {
		c.KeyFunc(c)
	}
}

func (c *Chip8) ClearKey(key byte) {
	c.key[key] = 0
	c.KeyFunc(c)
}

func (c *Chip8) PressKey(key byte) {
	for i := range c.key {
		if byte(i) == key {
			if c.key[i] != 1 {
				c.key[i] = 1
			}

		} else {
			if c.key[i] != 0 {
				c.key[i] = 0
			}
		}
	}
	c.KeyFunc(c)
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

func (c *Chip8) emulateCycle() {
	c.fetch()
	c.decode()
	c.updateTimers()
	if c.drawFlag {
		c.ScreenFunc(c)
		c.drawFlag = false
	}
	c.InfoFunc(c)
}

func (c *Chip8) run() {
	clock := time.NewTicker(clockCycleRate)
	timers := time.NewTicker(timeCycleRate)
	keyboard := time.NewTicker(keyboardCycleRate)

	go c.runClockCycle(clock)
	go c.runTimerCycle(timers)
	go c.runKeyboardCycle(keyboard)
}

func (c *Chip8) runKeyboardCycle(keyboardTimer *time.Ticker) {
	for {
		select {
		case <-c.stopSignal:
			keyboardTimer.Stop()
			return
		case <-keyboardTimer.C:
			c.ClearKeys()
		}
	}
}

func (c *Chip8) runClockCycle(clockTimer *time.Ticker) {
	for {
		select {
		case <-c.stopSignal:
			clockTimer.Stop()
			return
		case <-clockTimer.C:
			c.emulateCycle()

		case k := <-c.keyBoardInterrupt:
			c.PressKey(k)
		}
	}
}

func (c *Chip8) runTimerCycle(timerTimer *time.Ticker) {
	for {
		select {
		case <-c.stopSignal:
			timerTimer.Stop()
			return
		case <-timerTimer.C:
			c.updateTimers()
		}
	}
}

func (c *Chip8) stop() {
	c.stopSignal <- true
	c.stopSignal <- true
	c.stopSignal <- true
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

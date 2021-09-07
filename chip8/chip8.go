package chip8

import (
	"chip8/emulator"
	"io/ioutil"
	"time"
)

const (
	memorySize            = 4096
	vRegSize              = 16
	stackSize             = 16
	screenWidth           = 64
	screenHeigth          = 32
	keyNumbers            = 16
	clockCycleRate        = 2 * time.Microsecond
	timeCycleRate         = 16 * time.Microsecond
	keyboardResetDuration = 50 * time.Millisecond
)

var (
	keyboardInterrupt = make(chan byte, 1)
	stopSignal        = make(chan struct{})
	drawSignal        = make(chan []byte)
	keySignal         = make(chan []byte, 1)
	running           = false
)

type Chip8 struct {
	opcode       uint16
	i            uint16 // The address register, which is named I, is 12 bits wide and is used with several opcodes that involve memory operations.
	pc           uint16
	stack        [stackSize]uint16
	sp           byte
	memory       [memorySize]byte
	v            [vRegSize]byte // general purpose registers
	vChanged     [vRegSize]bool
	screenBuf    [screenWidth * screenHeigth]byte
	drawFlag     bool
	playingSound bool
	key          [keyNumbers]byte
	delayTimer   byte
	soundTimer   byte

	info       emulator.EmulatorInfo
	SetEmuInfo func(emulator.ChipGetter)
}

func (c *Chip8) Init(file string, tui emulator.TUISetter) error {
	c.pc = 0x200 // programs written for the original system begin at memory location 512 (0x200)
	c.SetEmuInfo = tui.SetEmuInfo
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

func (c Chip8) KeySignal() <-chan []byte {
	return keySignal
}

func (c Chip8) DrawSignal() <-chan []byte {
	return drawSignal
}

func (c *Chip8) fetch() {
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
}

func (c *Chip8) updateTimers() {
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.playingSound = true
		c.soundTimer--
	} else {
		c.playingSound = false
	}
}

func (c *Chip8) emulateCycle() {
	c.fetch()
	c.decode()
	c.updateTimers()
	if c.drawFlag {
		drawSignal <- c.screenBuf[:]
		c.drawFlag = false
	}
	c.SetEmuInfo(c)
}

func (c *Chip8) run() {
	stopSignal = make(chan struct{})
	clock := time.NewTicker(clockCycleRate)
	timers := time.NewTicker(timeCycleRate)
	running = true

	go c.runClockCycle(clock)
	go c.runTimerCycle(timers)
}

func (c *Chip8) runClockCycle(clockTimer *time.Ticker) {
	keyboardResetTimer := time.NewTimer(keyboardResetDuration)
	for {
		select {
		case <-stopSignal:
			clockTimer.Stop()
			return
		case <-clockTimer.C:
			c.emulateCycle()

		case k := <-keyboardInterrupt:
			keyboardResetTimer.Stop()
			keyboardResetTimer.Reset(keyboardResetDuration)
			c.pressKey(k)
		case <-keyboardResetTimer.C:
			c.clearKeys()
		}
	}
}

func (c *Chip8) runTimerCycle(timerTimer *time.Ticker) {
	for {
		select {
		case <-stopSignal:
			timerTimer.Stop()
			return
		case <-timerTimer.C:
			c.updateTimers()
		}
	}
}

func (c *Chip8) stop() {
	if running {
		close(stopSignal)
		running = false
	}
}

func (c *Chip8) setEmulatorInfo(n string, t string, d string) {
	c.info = emulator.CreateEmulatorInfo(c.opcode, n, t, d, c.pc)
}

func (c *Chip8) clearKeys() {
	clearedKey := false
	for i := range c.key {
		if c.key[i] != 0 {
			c.key[i] = 0
			clearedKey = true
		}
	}
	if clearedKey {
		keySignal <- c.key[:]
	}
}

func (c *Chip8) pressKey(key byte) {
	for i := range c.key {
		if byte(i) == key {
			if c.key[i] != 1 {
				c.key[i] = 1
			}
		}
	}
	keySignal <- c.key[:]
}

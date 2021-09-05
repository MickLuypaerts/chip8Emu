package view

import (
	"chip8/emulator"
	"fmt"
	"image"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	xOffSet       = 2
	yOffSet       = 4
	lMemRowLength = 16
)

type TUI struct {
	lGPR         *widgets.List
	lKeys        *widgets.List
	lStack       *widgets.List
	lMem         *widgets.List
	lProgStats   *widgets.List
	canvas       *ui.Canvas
	screenWidth  int
	screenHeight int
	grid         *ui.Grid
	termWidth    int
	termHeight   int
}

func (t *TUI) Init(c emulator.Chip) {
	t.initLGPR(c.GetGPRValues)
	t.initLKeys(c.GetKeyValues)
	t.initLStack(c.GetStackValues)
	t.initLMem(c)
	t.initLProgStats(c.OpcodeInfo)
	t.initCanvas()
	t.initTermSize()
	t.screenWidth, t.screenHeight = c.GetScreenSize()
	t.initGrid()
}

func (t *TUI) initLGPR(getGPRValues func() []string) {
	t.lGPR = widgets.NewList()
	t.lGPR.Title = "Registers"
	t.lGPR.Rows = getGPRValues()
	t.lGPR.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.lGPR.WrapText = false
	t.lGPR.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
}

func (t *TUI) initLKeys(getKeyValues func() []string) {
	t.lKeys = widgets.NewList()
	t.lKeys.Title = "Keys"
	t.lKeys.Rows = getKeyValues()
	t.lKeys.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.lKeys.WrapText = false
	t.lKeys.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
}

func (t *TUI) initLStack(getStackValues func() []string) {
	t.lStack = widgets.NewList()
	t.lStack.Title = "Stack"
	t.lStack.Rows = getStackValues()
	t.lStack.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.lStack.WrapText = false
	t.lStack.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
}
func (t *TUI) initLMem(c emulator.Chip) {
	t.lMem = widgets.NewList()
	t.lMem.Title = "Memory"
	t.lMem.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.lMem.WrapText = false

	t.lMem.Rows = memoryToTUIMemory(c.GetMemoryValues())
	t.SetListMemRow(c.OpcodeInfo())
}

func memoryToTUIMemory(memory []byte) []string {
	var mem []string
	var row string

	rowCount := 0
	columnCount := 0x0

	for i := 0; i < len(memory); i++ {
		if rowCount == 0 {
			row = fmt.Sprintf("0x%04X ", columnCount<<4)
			columnCount++
		}
		row += fmt.Sprintf("%02X ", memory[i])
		rowCount++
		if rowCount == lMemRowLength {
			mem = append(mem, row)
			rowCount = 0
		}
	}
	return mem
}
func (t *TUI) initLProgStats(getProgStats func() emulator.OpcodeInfo) {
	t.lProgStats = widgets.NewList()
	t.lProgStats.Title = fmt.Sprintf("INFO %s", os.Args[1])
	t.lProgStats.WrapText = true
	t.lProgStats.Rows = []string{fmt.Sprint(getProgStats())}
}
func (t *TUI) initCanvas() {
	t.canvas = ui.NewCanvas()
}

func (t *TUI) initGrid() {
	t.grid = ui.NewGrid()
	t.grid.SetRect(0, 0, t.termWidth, t.termHeight)
	t.grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(3.0/4, t.canvas),
			ui.NewCol(1.0/4, t.lProgStats),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(0.5/4, t.lGPR),
			ui.NewCol(0.5/4, t.lStack),
			ui.NewCol(0.5/4, t.lKeys),
			ui.NewCol(2.5/4, t.lMem),
		),
	)

}

func (t *TUI) initTermSize() {
	t.termWidth, t.termHeight = ui.TerminalDimensions()
}

func (t *TUI) SetEmuInfo(c emulator.ChipGetter) {
	t.lProgStats.Rows = []string{fmt.Sprint(c.OpcodeInfo())}
	t.lGPR.Rows = c.GetGPRValues()
	t.lMem.Rows = memoryToTUIMemory(c.GetMemoryValues())
	t.SetListMemRow(c.OpcodeInfo())
	t.lStack.Rows = c.GetStackValues()

	ui.Render(t.lProgStats, t.lGPR, t.lMem, t.lStack)
}

func (t *TUI) SetListMemRow(opcodeInfo emulator.OpcodeInfo) {
	row := int(opcodeInfo.ProgramCount() / lMemRowLength)
	if row > len(t.lMem.Rows)-1 {
		t.lMem.SelectedRow = len(t.lMem.Rows) - 1
	} else if row < 0 {
		t.lMem.SelectedRow = 0
	} else {
		t.lMem.SelectedRow = row
	}
	ui.Render(t.lMem)
}

func (t *TUI) SetKeyInfo(c emulator.ChipGetter) {
	t.lKeys.Rows = c.GetKeyValues()
	ui.Render(t.lKeys)
}

func (t *TUI) UpdateScreen(c emulator.ChipGetter) {
	screenBuffer := c.GetScreen()
	for y := 0; y < t.screenHeight; y++ {
		for x := 0; x < t.screenWidth; x++ {
			if screenBuffer[x+(y*t.screenWidth)] != 0 {
				t.canvas.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorRed)
			} else {
				t.canvas.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorWhite)
			}
		}
	}
	ui.Render(t.canvas)
}

func scrollDown(l *widgets.List) {
	l.ScrollDown()
	ui.Render(l)
}

func scrollUp(l *widgets.List) {
	l.ScrollUp()
	ui.Render(l)
}

func scrollTop(l *widgets.List) {
	l.ScrollTop()
	ui.Render(l)
}

func scrollBottom(l *widgets.List) {
	l.ScrollBottom()
	ui.Render(l)
}

func (t TUI) Close() {
	ui.Close()
}

func (t TUI) Render() {
	ui.Render(t.grid)
}

func (t TUI) Setup() error {
	if err := ui.Init(); err != nil {
		return err
	}
	return nil
}

func (t TUI) KeyEvent() <-chan string {
	ch := make(chan string)
	keyEvents := ui.PollEvents()
	go func() {
		for {
			key := <-keyEvents
			ch <- key.ID
		}
	}()
	return ch
}

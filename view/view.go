package view

import (
	"fmt"
	"image"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	xOffSet = 2
	yOffSet = 4
)

type TUI struct {
	lGPR       *widgets.List
	lKeys      *widgets.List
	lStack     *widgets.List
	LMem       *widgets.List
	lProgStats *widgets.List
	canvas     *ui.Canvas
	Grid       *ui.Grid
	termWidth  int
	termHeight int
}

type Chip interface {
	GetGPRValues() []string
	GetKeyValues() []string
	GetStackValues() []string
	GetMemoryValues() []string
	GetMemoryRow() uint16
	GetProgStats() []string
	GetScreen() ([]byte, int, int)
}

func (t *TUI) Init(c Chip) {
	t.initLGPR(c.GetGPRValues)
	t.initLKeys(c.GetKeyValues)
	t.initLStack(c.GetStackValues)
	t.initLMem(c.GetMemoryValues, c.GetMemoryRow)
	t.initLProgStats(c.GetProgStats)
	t.initCanvas()
	t.initTermSize()
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
func (t *TUI) initLMem(getMemoryValues func() []string, getMemoryRow func() uint16) {
	t.LMem = widgets.NewList()
	t.LMem.Title = "Memory"
	t.LMem.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.LMem.WrapText = false
	t.LMem.Rows = getMemoryValues()
	t.LMem.SelectedRow = int(getMemoryRow())
}
func (t *TUI) initLProgStats(getProgStats func() []string) {
	t.lProgStats = widgets.NewList()
	t.lProgStats.Title = fmt.Sprintf("INFO %s", os.Args[1])
	t.lProgStats.WrapText = true
	t.lProgStats.Rows = getProgStats()
}
func (t *TUI) initCanvas() {
	t.canvas = ui.NewCanvas()
}

func (t *TUI) initGrid() {
	t.Grid = ui.NewGrid()
	t.Grid.SetRect(0, 0, t.termWidth, t.termHeight)
	t.Grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(3.0/4, t.canvas),
			ui.NewCol(1.0/4, t.lProgStats),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(0.5/4, t.lGPR),
			ui.NewCol(0.5/4, t.lStack),
			ui.NewCol(0.5/4, t.lKeys),
			ui.NewCol(2.5/4, t.LMem),
		),
	)

}

func (t *TUI) initTermSize() {
	t.termWidth, t.termHeight = ui.TerminalDimensions()
}

func (t *TUI) SetEmuInfo(c Chip) {
	t.lProgStats.Rows = c.GetProgStats()
	t.lGPR.Rows = c.GetGPRValues()
	t.LMem.Rows = c.GetMemoryValues()
	t.LMem.SelectedRow = int(c.GetMemoryRow())
	t.lStack.Rows = c.GetStackValues()
	ui.Render(t.lProgStats, t.lGPR, t.LMem, t.lStack)
}

func (t *TUI) UpdateScreen(c Chip) {
	screenBuffer, width, height := c.GetScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if screenBuffer[x+(y*width)] != 0 {
				t.canvas.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorRed)
			} else {
				t.canvas.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorWhite)
			}
		}
	}
	ui.Render(t.canvas)
}

func ScrollDown(l *widgets.List) {
	l.ScrollDown()
	ui.Render(l)
}

func ScrollUp(l *widgets.List) {
	l.ScrollUp()
	ui.Render(l)
}

func ScrollTop(l *widgets.List) {
	l.ScrollTop()
	ui.Render(l)
}

func ScrollBottom(l *widgets.List) {
	l.ScrollBottom()
	ui.Render(l)
}

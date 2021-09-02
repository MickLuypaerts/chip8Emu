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
	LGPR       *widgets.List
	LKeys      *widgets.List
	LStack     *widgets.List
	LMem       *widgets.List
	LProgStats *widgets.List
	Canvas     *ui.Canvas
	Grid       *ui.Grid
	TermWidth  int
	TermHeight int
}

type chip interface {
	GetGPRValues() []string
	GetKeyValues() []string
	GetStackValues() []string
	GetMemoryValues() []string
	GetMemoryRow() uint16
	GetProgStats() []string
	GetScreen() ([]byte, int, int)
}

func (t *TUI) Init(c chip) {
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
	t.LGPR = widgets.NewList()
	t.LGPR.Title = "Registers"
	t.LGPR.Rows = getGPRValues()
	t.LGPR.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.LGPR.WrapText = false
	t.LGPR.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
}

func (t *TUI) initLKeys(getKeyValues func() []string) {
	t.LKeys = widgets.NewList()
	t.LKeys.Title = "Keys"
	t.LKeys.Rows = getKeyValues()
	t.LKeys.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.LKeys.WrapText = false
	t.LKeys.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
}

func (t *TUI) initLStack(getStackValues func() []string) {
	t.LStack = widgets.NewList()
	t.LStack.Title = "Stack"
	t.LStack.Rows = getStackValues()
	t.LStack.TextStyle = ui.NewStyle(ui.ColorYellow)
	t.LStack.WrapText = false
	t.LStack.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
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
	t.LProgStats = widgets.NewList()
	t.LProgStats.Title = fmt.Sprintf("INFO %s", os.Args[1])
	t.LProgStats.WrapText = true
	t.LProgStats.Rows = getProgStats()
}
func (t *TUI) initCanvas() {
	t.Canvas = ui.NewCanvas()
}

func (t *TUI) initGrid() {
	t.Grid = ui.NewGrid()
	t.Grid.SetRect(0, 0, t.TermWidth, t.TermHeight)
	t.Grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(3.0/4, t.Canvas),
			ui.NewCol(1.0/4, t.LProgStats),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(0.5/4, t.LGPR),
			ui.NewCol(0.5/4, t.LStack),
			ui.NewCol(0.5/4, t.LKeys),
			ui.NewCol(2.5/4, t.LMem),
		),
	)

}

func (t *TUI) initTermSize() {
	t.TermWidth, t.TermHeight = ui.TerminalDimensions()
}

func (t *TUI) SetEmuInfo(c chip) {
	t.LProgStats.Rows = c.GetProgStats()
	t.LGPR.Rows = c.GetGPRValues()
	t.LMem.Rows = c.GetMemoryValues()
	t.LMem.SelectedRow = int(c.GetMemoryRow())
	t.LStack.Rows = c.GetStackValues()
	ui.Render(t.LProgStats, t.LGPR, t.LMem, t.LStack)
}

func (t *TUI) UpdateScreen(c chip) {
	screenBuffer, width, height := c.GetScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if screenBuffer[x+(y*width)] != 0 {
				t.Canvas.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorRed)
			} else {
				t.Canvas.SetPoint(image.Pt(x*xOffSet, y*yOffSet), ui.ColorWhite)
			}
		}
	}
	ui.Render(t.Canvas)
}

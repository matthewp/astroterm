package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BuildUI struct {
	Flex *tview.Flex
}

func NewBuildUI() *BuildUI {
	flex := tview.NewFlex()

	overview := tview.NewFlex()
	overview.SetBorder(true)
	overview.SetTitle("Overview")

	flex.AddItem(overview, 0, 1, false)

	b := &BuildUI{
		Flex: flex,
	}

	return b
}

/* UISection implementation */
func (b *BuildUI) Primitive() tview.Primitive {
	return b.Flex
}

func (b *BuildUI) Stop() {}

func (b *BuildUI) SetFocusMenu(focusMenu func()) {

}

func (b *BuildUI) MakeActive(cmds *BottomCommandsUI) {

}

func (b *BuildUI) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	return event
}

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type IntegrationsUI struct {
	Flex *tview.Flex
}

func NewIntegrationsUI() *IntegrationsUI {
	flex := tview.NewFlex()

	overview := tview.NewFlex()
	overview.SetBorder(true)
	overview.SetTitle("Integrations")

	flex.AddItem(overview, 0, 1, false)

	it := &IntegrationsUI{
		Flex: flex,
	}

	return it
}

/* UISection implementation */
func (b *IntegrationsUI) Primitive() tview.Primitive {
	return b.Flex
}

func (b *IntegrationsUI) Stop() {}

func (b *IntegrationsUI) SetFocusMenu(focusMenu func()) {

}

func (b *IntegrationsUI) MakeActive(cmds *BottomCommandsUI) {

}

func (b *IntegrationsUI) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	return event
}

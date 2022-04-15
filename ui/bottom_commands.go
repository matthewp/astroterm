package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

type BottomCommandsUI struct {
	*tview.Flex
	app     *tview.Application
	navForm *tview.Form
}

func NewBottomCommands(app *tview.Application) *BottomCommandsUI {
	flex := tview.NewFlex()

	navForm := tview.NewForm()
	navForm.SetBorderPadding(0, 0, 0, 0)
	navForm.SetBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonTextColor(NavStyles.TextColor)

	flex.AddItem(navForm, 0, 1, false)

	return &BottomCommandsUI{
		Flex:    flex,
		app:     app,
		navForm: navForm,
	}
}

func (c *BottomCommandsUI) FormatLabel(label string, shortcut rune) string {
	return fmt.Sprintf("%s [[#be0000::b]%c[-:-:-]]", label, shortcut)
}

func (c *BottomCommandsUI) AddButton(label string, shortcut rune, cb func()) *tview.Button {
	idx := c.navForm.GetButtonCount()
	lbl := c.FormatLabel(label, shortcut)
	c.navForm.AddButton(lbl, cb)
	btn := c.navForm.GetButton(idx)
	return btn
}

func (c *BottomCommandsUI) ClearButtons() {
	c.navForm.ClearButtons()
}

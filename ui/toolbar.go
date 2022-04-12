package ui

import (
	"github.com/rivo/tview"
)

type Toolbar struct {
	*tview.Flex
	app *tview.Application
}

func NewToolbar(app *tview.Application) *Toolbar {
	flex := tview.NewFlex()

	navForm := tview.NewForm()
	navForm.SetBorderPadding(0, 0, 0, 0)
	navForm.SetBackgroundColor(NavStyles.BackgroundColor)
	navForm.AddButton("[#be0000::b]F[-:-:-]ile", nil)
	navForm.AddButton("[#be0000::b]W[-:-:-]orkspaces", nil)
	navForm.AddButton("[#be0000::b]H[-:-:-]elp", nil)
	navForm.SetButtonBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonTextColor(NavStyles.TextColor)

	flex.AddItem(navForm, 0, 1, false)

	return &Toolbar{
		Flex: flex,
		app:  app,
	}
}

package ui

import (
	"github.com/rivo/tview"
)

type MainNav struct {
	*tview.Flex
	app *tview.Application
}

func NewMainNav(app *tview.Application) *MainNav {
	flex := tview.NewFlex()

	navForm := tview.NewForm()
	navForm.SetBorderPadding(0, 0, 0, 0)
	navForm.SetBackgroundColor(NavStyles.BackgroundColor)
	navForm.AddButton("[#be0000::b]D[-:-:-]evelopment", nil)
	navForm.AddButton("[#be0000::b]P[-:-:-]roduction", nil)
	navForm.SetButtonBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonTextColor(NavStyles.TextColor)

	flex.AddItem(navForm, 0, 1, false)

	return &MainNav{
		Flex: flex,
		app:  app,
	}
}

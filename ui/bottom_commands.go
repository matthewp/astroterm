package ui

import "github.com/rivo/tview"

type BottomCommands struct {
	*tview.Flex
	app *tview.Application
}

func NewBottomCommands(app *tview.Application) *BottomCommands {
	flex := tview.NewFlex()

	navForm := tview.NewForm()
	navForm.SetBorderPadding(0, 0, 0, 0)
	navForm.SetBackgroundColor(NavStyles.BackgroundColor)
	navForm.AddButton("One", nil)
	navForm.AddButton("Two", nil)
	navForm.SetButtonBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonTextColor(NavStyles.TextColor)

	flex.AddItem(navForm, 0, 1, false)

	return &BottomCommands{
		Flex: flex,
		app:  app,
	}
}

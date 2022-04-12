package ui

import "github.com/rivo/tview"

type MainNav struct {
	*tview.Flex
	app *tview.Application
}

func NewMainNav(app *tview.Application) *MainNav {
	flex := tview.NewFlex()
	/*flex.SetTitle("Nav")
	flex.SetTitleAlign(tview.AlignLeft)
	flex.SetBorder(true)
	flex.SetDirection(tview.FlexRow)*/

	navForm := tview.NewForm()
	navForm.SetBorderPadding(0, 0, 0, 0)
	navForm.SetBackgroundColor(NavStyles.BackgroundColor)
	navForm.AddButton("[red::b]D[-][black::b]evelopment", nil)
	navForm.AddButton("[red::b]P[-][::b]roduction", nil)
	navForm.SetButtonBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonTextColor(NavStyles.TextColor)

	flex.AddItem(navForm, 0, 1, false)

	dForm := tview.NewForm()
	dForm.SetBorderPadding(0, 0, 0, 0)
	//dForm.SetTitle("Server")
	dForm.SetBorder(false)
	dForm.SetButtonsAlign(tview.AlignRight)

	dForm.AddButton("Start dev server", nil)
	dForm.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	flex.AddItem(dForm, 25, 0, false)

	//flex.AddItem(btn, 0, 1, false)

	return &MainNav{
		Flex: flex,
		app:  app,
	}
}

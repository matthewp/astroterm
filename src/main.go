package main

import (
	"astroterm/src/ui"

	"github.com/rivo/tview"
)

func main() {
	start()
}

func testing() {

	app := tview.NewApplication().EnableMouse(true)
	button := tview.NewButton("Hit Enter to close").SetSelectedFunc(func() {
		app.Stop()
	})
	button.SetBorder(true).SetRect(0, 0, 22, 3)
	/*button.SetBackgroundColorActivated(tcell.ColorRebeccaPurple)
	button.SetBackgroundColor(tcell.ColorAliceBlue)
	button.SetLabelColor(tcell.ColorRebeccaPurple)*/
	if err := app.SetRoot(button, false).SetFocus(button).Run(); err != nil {
		panic(err)
	}

}

func start() {
	u := ui.NewUI()
	if err := u.Start(); err != nil {
		panic(err)
	}
}

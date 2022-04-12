package ui

import (
	"github.com/rivo/tview"
)

type FixedControls struct {
	*tview.Grid
}

func NewFixedControls() *FixedControls {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Process")

	/*dForm := tview.NewForm()
	dForm.SetBorderPadding(0, 0, 0, 0)
	dForm.SetBorder(false)
	dForm.SetButtonsAlign(tview.AlignRight)

	dForm.AddButton("Start dev server", nil)
	dForm.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)*/
	form.AddButton("Start dev server", nil)
	form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	form.SetButtonsAlign(tview.AlignCenter)

	grid := tview.NewGrid().
		SetColumns(0, 30, 1).
		SetRows(1, 5, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	return &FixedControls{
		Grid: grid,
	}
}

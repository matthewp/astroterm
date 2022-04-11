package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app  *tview.Application
	grid *tview.Grid
}

func NewUI() *UI {
	newPrimitive := func(text string) tview.Primitive {
		tv := tview.NewTextView()

		tv.SetTitle(text)
		tv.SetTitleAlign(tview.AlignLeft)
		tv.SetBorder(true)

		return tv
	}

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
			return nil
		}
		return event
	})

	menu := newPrimitive("Menu")
	main := NewDevServer(app)
	sideBar := newPrimitive("Side Bar")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		//SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	return &UI{
		app:  app,
		grid: grid,
	}
}

func (ui *UI) Start() error {
	app := ui.app
	grid := ui.grid
	app.SetRoot(grid, true).SetFocus(grid).EnableMouse(true)

	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func otherStuff() {
	box := tview.NewBox().SetBorder(true).SetTitle("Commands")

	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

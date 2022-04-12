package ui

import (
	aenv "astroterm/env"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app   *tview.Application
	grid  *tview.Grid
	pages *tview.Pages
}

func NewUI() *UI {
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 /* q */ {
			app.Stop()
			return nil
		}
		return event
	})

	grid := newGrid(app)
	pages := tview.NewPages().
		AddPage("background", grid, true, true).
		AddPage("modal", NewFixedControls(), true, true)

	return &UI{
		app:   app,
		grid:  grid,
		pages: pages,
	}
}

func newGrid(app *tview.Application) *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		tv := tview.NewTextView()

		tv.SetTitle(text)
		tv.SetTitleAlign(tview.AlignLeft)
		tv.SetBorder(true)

		return tv
	}

	menu := newPrimitive("Menu")
	nav := NewMainNav(app)
	main := NewDevServer(app)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(30, 0, 30).
		AddItem(nav, 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 2, 0, 100, false)

	return grid
}

func (ui *UI) Start() error {
	app := ui.app

	env, err := aenv.GetEnvironment()

	if env == nil {
		return err
	}

	if env.IsAstroProject {
		app.SetRoot(ui.pages, true).SetFocus(ui.grid).EnableMouse(true)
	} else {
		naModal := notAnAstroAppModal(app)
		app.SetRoot(naModal, true).SetFocus(naModal).EnableMouse(true)
	}

	if err = app.Run(); err != nil {
		return err
	}
	return nil
}

func notAnAstroAppModal(app *tview.Application) *tview.Modal {
	modal := tview.NewModal().
		SetText("This does not appear to be an Astro project. Change into a directory that contains an Astro project and start astroterm again.").
		AddButtons([]string{"Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			}
		})
	return modal
}

func otherStuff() {
	box := tview.NewBox().SetBorder(true).SetTitle("Commands")

	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

package ui

import (
	aenv "astroterm/env"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app         *tview.Application
	grid        *tview.Grid
	menu        *Menu
	main        *tview.Flex
	currentMain UISection
	pages       *tview.Pages
}

type UISectionType int64

const (
	SectionDevelopment UISectionType = iota
	SectionBuild
	SectionIntegrations
	SectionDiagnostics
)

type UISection interface {
	Primitive() tview.Primitive
	Stop()
}

func NewUI() *UI {
	app := tview.NewApplication()
	ui := &UI{
		app: app,
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 /* q */ {
			if ui.currentMain != nil {
				ui.currentMain.Stop()
			}

			app.Stop()
			return nil
		}
		return event
	})

	toolbar := NewToolbar(app)
	menu := NewMenu(ui)
	main := tview.NewFlex()
	ui.menu = menu
	ui.main = main

	defaultMain := NewDevServer(app)
	ui.SetMainItem(defaultMain)

	cmds := NewBottomCommands(app)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(30, 0, 30).
		AddItem(toolbar, 0, 0, 1, 3, 0, 0, false).
		AddItem(cmds, 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 2, 0, 100, false)

	pages := tview.NewPages().AddPage("background", grid, true, true)

	ui.pages = pages
	ui.grid = grid

	return ui
}

func (ui *UI) Start() error {
	app := ui.app

	env, err := aenv.GetEnvironment()

	if env == nil {
		return err
	}

	if env.IsAstroProject {
		app.SetRoot(ui.pages, true).SetFocus(ui.menu).EnableMouse(true)
	} else {
		naModal := notAnAstroAppModal(app)
		app.SetRoot(naModal, true).SetFocus(naModal).EnableMouse(true)
	}

	if err = app.Run(); err != nil {
		return err
	}
	return nil
}

func (u *UI) Navigate(sec UISectionType) {
	switch sec {
	case SectionDevelopment:
		u.main.RemoveItem(u.currentMain.Primitive())
		u.SetMainItem(NewDevServer(u.app))
		break
	case SectionBuild:
		u.main.RemoveItem(u.currentMain.Primitive())
		u.SetMainItem(NewBuildUI())
		break
	case SectionIntegrations:
		u.main.RemoveItem(u.currentMain.Primitive())
		u.SetMainItem(NewIntegrationsUI())
		break
	case SectionDiagnostics:
		u.main.RemoveItem(u.currentMain.Primitive())
		break
	}
}

func (u *UI) SetMainItem(item UISection) {
	p := item.Primitive()
	u.main.AddItem(p, 0, 1, false)
	u.currentMain = item
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

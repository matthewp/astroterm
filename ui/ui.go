package ui

import (
	"astroterm/db"
	aenv "astroterm/env"
	"astroterm/project"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	DevModel       *db.DevServerModel
	CurrentProject *project.Project

	db          *db.Database
	app         *tview.Application
	grid        *tview.Grid
	menu        *Menu
	main        *tview.Flex
	currentMain UISection
	pages       *tview.Pages

	sections map[UISectionType]UISection
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
		DevModel:       &db.DevServerModel{},
		CurrentProject: loadLocalProject(),
		app:            app,
		db:             db.NewDatabase(),
		sections:       make(map[UISectionType]UISection),
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

	defaultMain := ui.LoadSection(SectionDevelopment)
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
	newSec := u.LoadSection(sec)
	u.main.RemoveItem(u.currentMain.Primitive())
	u.SetMainItem(newSec)
}

func (u *UI) LoadSection(sec UISectionType) UISection {
	sections := u.sections
	if val, ok := sections[sec]; ok {
		return val
	}
	var val UISection
	switch sec {
	case SectionDevelopment:
		val = NewDevServer(u)
		break
	case SectionBuild:
		val = NewBuildUI()
		break
	case SectionIntegrations:
		val = NewIntegrationsUI()
		break
	case SectionDiagnostics:
		val = nil
		break
	default:
		val = nil
	}
	if val != nil {
		sections[sec] = val
	}
	return val
}

func (u *UI) SetMainItem(item UISection) {
	p := item.Primitive()
	u.main.AddItem(p, 0, 1, false)
	u.currentMain = item
}

func (u *UI) Draw() *tview.Application {
	return u.app.Draw()
}

func (u *UI) SetFocus(p tview.Primitive) *tview.Application {
	return u.app.SetFocus(p)
}

func loadLocalProject() *project.Project {
	var localProject *project.Project
	localProject, err := project.OpenLocalProject()
	if err != nil {
		localProject = project.NewProject()
	}
	return localProject
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

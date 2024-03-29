package ui

import (
	"astroterm/actors"
	"astroterm/db"
	aenv "astroterm/env"
	"astroterm/project"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	DevModel       *db.DevServerModel
	CurrentProject *project.Project
	DB             *db.Database
	app            *tview.Application
	DevActor       *actors.DevServerActor
	BuildActor     *actors.BuildActor
	grid           *tview.Grid
	menu           *Menu
	cmds           *BottomCommandsUI
	main           *tview.Flex
	currentMain    UISection
	pages          *tview.Pages
	scBtn          *tview.Button
	sections       map[UISectionType]UISection
	toolbar        *Toolbar
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
	MakeActive(*BottomCommandsUI)
	Stop()
	SetFocusMenu(func())
	InputCapture(event *tcell.EventKey) *tcell.EventKey
}

func NewUI() *UI {
	app := tview.NewApplication()
	ui := &UI{
		DevModel:       &db.DevServerModel{},
		CurrentProject: loadLocalProject(),
		app:            app,
		DB:             db.NewDatabase(),
		sections:       make(map[UISectionType]UISection),
	}
	ui.DevActor = actors.NewDevServerActor(ui.CurrentProject).Start()
	ui.BuildActor = actors.NewBuildActor(ui.CurrentProject).Start()
	ui.toolbar = NewToolbar(ui, ui.DevActor, ui.BuildActor)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			ui.MaybeStop()
			return nil
		}
		if ui.toolbar.HasOpenMenu() {
			return ui.toolbar.InputCapture(event)
		} else {
			event = ui.currentMain.InputCapture(event)
			if event != nil {
				ui.toolbar.InputCapture(event)
			}
			return event
		}
	})

	ui.toolbar.SetProject(ui.CurrentProject)
	menu := NewMenu(ui)
	menu.SetFocusSection(func() {
		if ui.currentMain != nil {
			ui.app.SetFocus(ui.currentMain.Primitive())
		}
	})
	main := tview.NewFlex()
	ui.menu = menu
	ui.main = main

	cmds := NewBottomCommands(app)
	ui.cmds = cmds

	defaultMain := ui.LoadSection(SectionDevelopment)
	ui.SetMainItem(defaultMain)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(30, 0, 30).
		//AddItem(toolbar, 0, 0, 1, 3, 0, 0, false).
		AddItem(cmds, 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 2, 0, 100, false)

	pages := tview.NewPages()
	pages.AddPage("grid", grid, true, true)
	pages.AddPage("toolbar", ui.toolbar, true, true)
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

		//app.SetFocus(ui.toolbar.ml)
	} else {
		naModal := notAnAstroAppModal(app)
		app.SetRoot(naModal, true).SetFocus(naModal).EnableMouse(true)
	}

	if err = app.Run(); err != nil {
		return err
	}
	return nil
}

func (ui *UI) MaybeStop() {
	if ui.ConfirmStop() {
		ui.app.Stop()
	}
}

func (ui *UI) ConfirmStop() bool {
	var devServerRunning bool = ui.DevModel.Pid != 0
	if devServerRunning {
		app := ui.app
		pages := ui.pages
		modal := tview.NewModal().
			SetText("The [::b]dev server[-:-:-] is running. Would you like to stop it?").
			AddButtons([]string{"Stop server", "Just quit", "Cancel"}).
			SetDoneFunc(func(idx int, label string) {
				switch idx {
				case 0:
					ui.LoadSection(SectionDevelopment).Stop()
					app.Stop()
					break
				case 1:
					app.Stop()
				case 2:
					pages.RemovePage("modal")
					pages.SwitchToPage("grid")
				}
			})
		pages.AddPage("modal", modal, true, true)
		pages.SwitchToPage("modal")
		return false
	}
	return true
}

func (u *UI) Navigate(sec UISectionType) {
	newSec := u.LoadSection(sec)
	u.main.RemoveItem(u.currentMain.Primitive())
	if newSec != nil {
		u.SetMainItem(newSec)
	}

}

func (u *UI) LoadSection(sec UISectionType) UISection {
	sections := u.sections
	if val, ok := sections[sec]; ok {
		return val
	}
	var val UISection
	switch sec {
	case SectionDevelopment:
		val = NewDevServer(u, u.DevActor)
		break
	case SectionBuild:
		val = NewBuildUI(u, u.BuildActor)
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
		val.SetFocusMenu(u.focusMenu)
		sections[sec] = val
	}
	return val
}

func (u *UI) SetMainItem(item UISection) {
	p := item.Primitive()
	u.main.AddItem(p, 0, 1, false)
	u.cmds.ClearButtons()
	item.MakeActive(u.cmds)
	u.currentMain = item

}

func (u *UI) Draw() *tview.Application {
	return u.app.Draw()
}

func (u *UI) SetFocus(p tview.Primitive) *tview.Application {
	return u.app.SetFocus(p)
}

func (u *UI) GetFocus() tview.Primitive {
	return u.app.GetFocus()
}

func (u *UI) focusMenu() {
	u.app.SetFocus(u.menu)
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

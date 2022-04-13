package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Menu struct {
	*tview.List
	ui *UI
}

func NewMenu(u *UI) *Menu {
	list := tview.NewList().
		AddItem("Development", "", 'd', nil).
		AddItem("Build", "", 'b', nil).
		AddItem("Integrations", "", 'i', nil).
		AddItem("Diagnostics", "", 'n', nil)
		//AddItem("List item 4", "Some explanatory text", 'd', nil)
	list.SetBorder(true)
	list.SetTitle("Menu")
	list.SetTitleAlign(tview.AlignLeft)
	list.ShowSecondaryText(false)
	list.SetHighlightFullLine(true)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyTab:
			return nil
		case tcell.KeyBacktab:
			return nil
		}
		return event
	})

	list.SetChangedFunc(func(idx int, mainText string, secondaryText string, shortcut rune) {
		switch idx {
		case 0:
			u.Navigate(SectionDevelopment)
			break
		case 1:
			u.Navigate(SectionBuild)
			break
		case 2:
			u.Navigate(SectionIntegrations)
			break
		case 3:
			u.Navigate(SectionDiagnostics)
			break
		}
	})

	return &Menu{
		List: list,
		ui:   u,
	}
}

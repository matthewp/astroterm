package ui

import (
	"astroterm/project"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Toolbar struct {
	*tview.Flex
	app      *tview.Application
	titlebar *tview.TextView
}

var titleColors = []tcell.Color{
	tcell.NewRGBColor(199, 186, 185),
	tcell.NewRGBColor(29, 34, 40),
	tcell.NewRGBColor(114, 65, 65),
	tcell.NewRGBColor(118, 168, 176),
	tcell.NewRGBColor(102, 133, 159),
}

const pulsateDuration = 300

func NewToolbar(app *tview.Application) *Toolbar {
	flex := tview.NewFlex()

	navForm := tview.NewForm()
	navForm.SetBorderPadding(0, 0, 0, 0)
	navForm.SetBackgroundColor(NavStyles.BackgroundColor)
	navForm.AddButton("[#be0000::b]F[-:-:-]ile", nil)
	navForm.AddButton("[#be0000::b]W[-:-:-]orkspaces", nil)
	navForm.AddButton("[#be0000::b]H[-:-:-]elp", nil)
	navForm.SetButtonBackgroundColor(NavStyles.BackgroundColor)
	navForm.SetButtonTextColor(NavStyles.TextColor)

	titlebar := tview.NewTextView()
	titlebar.SetBackgroundColor(NavStyles.BackgroundColor)
	titlebar.SetTextColor(NavStyles.TextColor)
	titlebar.SetTextAlign(tview.AlignRight)
	titlebar.SetBorderPadding(0, 0, 0, 2)

	flex.AddItem(navForm, 0, 2, false)
	flex.AddItem(titlebar, 0, 1, false)

	t := &Toolbar{
		Flex:     flex,
		app:      app,
		titlebar: titlebar,
	}

	//go pulsateTitle(t, 0, true)
	return t
}

func pulsateTitle(t *Toolbar, idx int, forward bool) {
	color := titleColors[idx]
	t.titlebar.SetTextColor(color)
	t.app.Draw()

	time.Sleep(pulsateDuration * time.Millisecond)
	nf := forward
	nidx := idx
	if forward {
		lastIdx := len(titleColors) - 1
		if lastIdx == idx {
			nf = false
		} else {
			nidx = idx + 1
		}
	} else {
		if idx == 0 {
			nf = true
		} else {
			nidx = idx - 1
		}
	}
	pulsateTitle(t, nidx, nf)
}

func (t *Toolbar) SetProject(p *project.Project) {
	t.titlebar.SetText(p.Name())
}

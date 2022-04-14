package ui

import (
	"astroterm/project"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/tview"
)

type Toolbar struct {
	*tview.Flex
	app      *tview.Application
	titlebar *tview.TextView
}

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

	go pulsateTitle(t, getColors(), 0, true)
	return t
}

const pulsateDuration = 50

func getColors() []tcell.Color {
	c1, _ := colorful.Hex("#cc2b5e")
	c2, _ := colorful.Hex("#753a88")

	blocks := 50
	colors := make([]tcell.Color, blocks)

	for i := 0; i < blocks; i++ {
		c := c1.BlendRgb(c2, float64(i)/float64(blocks-1))
		r, g, b := c.RGB255()
		colors[i] = tcell.NewRGBColor(int32(r), int32(g), int32(b))
	}
	return colors
}

func pulsateTitle(t *Toolbar, titleColors []tcell.Color, idx int, forward bool) {
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
	pulsateTitle(t, titleColors, nidx, nf)
}

func (t *Toolbar) SetProject(p *project.Project) {
	t.titlebar.SetText(p.Name())
}

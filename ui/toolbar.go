package ui

import (
	"astroterm/project"
	"astroterm/util"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/matthewp/bestbar"
)

type Toolbar struct {
	*bestbar.Toolbar
	u       *UI
	servers int
}

func NewToolbar(u *UI) *Toolbar {
	t := bestbar.NewToolbar()
	t.SetDrawFunc(func() {
		u.Draw()
	})

	t.AddMenuList("File", 'F').
		AddItem("Open", 'O', nil).
		AddItem("Exit", 'x', func() {
			u.MaybeStop()
		})
	t.AddMenuList("🟢 Run", 'R').
		AddItem("Build", 'B', func() {

		}).
		AddItem("Start dev server", 'd', func() {

		}).
		AddItem("Start the preview server", 'p', func() {

		})
	t.AddMenuList("Help", 'H').
		AddItem("Documentation", 'D', func() {
			util.OpenBrowser("https://pkg.spooky.click/astroterm/")
		})

	tb := &Toolbar{
		Toolbar: t,
		u:       u,
		servers: 0,
	}

	go pulsateTitle(tb, getColors(), 0, true)

	return tb
}

func (t *Toolbar) SetDevServerRunning(running bool) {

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
	t.Toolbar.SetTitleTextColor(color)
	t.u.Draw()

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
	t.SetTitle(p.Name())
}

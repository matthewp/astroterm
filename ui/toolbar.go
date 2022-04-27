package ui

import (
	"astroterm/actors"
	"astroterm/project"
	"astroterm/util"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/matthewp/bestbar"
)

type toolbarState struct {
	devRunning     bool
	previewRunning bool
}

type Toolbar struct {
	*bestbar.Toolbar
	u               *UI
	state           *toolbarState
	servers         int
	devActor        *actors.DevServerActor
	runList         *bestbar.MenuList
	devServerAction *bestbar.MenuListItem
}

func NewToolbar(u *UI, devActor *actors.DevServerActor) *Toolbar {
	t := bestbar.NewToolbar()
	t.SetDrawFunc(func() {
		u.Draw()
	})

	tb := &Toolbar{
		state:    &toolbarState{},
		Toolbar:  t,
		u:        u,
		servers:  0,
		devActor: devActor,
	}

	t.AddMenuList("File", 'F').
		AddItem("Open", 'O', nil).
		AddItem("Exit", 'x', func() {
			u.MaybeStop()
		})
	tb.runList = t.AddMenuList("🔴 Run", 'R').
		AddItem("Build", 'B', func() {

		}).
		AddItem("Start dev server", 'd', func() {
			if tb.state.devRunning {
				tb.devActor.StopDevServer()
			} else {
				tb.devActor.StartDevServer()
			}
		}).
		AddItem("Start the preview server", 'p', func() {

		})
	tb.devServerAction = tb.runList.GetItem(1)
	t.AddMenuList("Help", 'H').
		AddItem("Documentation", 'D', func() {
			util.OpenBrowser("https://pkg.spooky.click/astroterm/")
		})

	// Event listeners
	go tb.listenToEvents()

	go pulsateTitle(tb, getColors(), 0, true)

	return tb
}

// View update functions
func (t *Toolbar) setRunButtonLabel(label string) {
	t.runList.SetButtonLabel(label, 'R')
}

func (t *Toolbar) setDevActionLabel(label string) {
	t.devServerAction.SetLabel(label, 'd')
}

// State update functions
func (t *Toolbar) setDevServerRunning(running bool) {
	if t.state.devRunning != running {
		t.state.devRunning = running

		if running {
			t.setDevActionLabel("Stop dev server")
			t.setRunButtonLabel("🟢 Run")
		} else {
			t.setDevActionLabel("Start dev server")
			t.setRunButtonLabel("🔴 Run")
		}
	}
}

// Event listeners
func (t *Toolbar) listenToEvents() {
	//ischan := t.devActor.SubscribeToInitialState()
	schan := t.devActor.SubscribeToStarting()
	stchan := t.devActor.SubscribeToStopped()
	for {
		select {
		case _ = <-schan:
			t.setDevServerRunning(true)
			break
		case _ = <-stchan:
			t.setDevServerRunning(false)
			break
		}
	}
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

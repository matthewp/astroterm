package ui

import (
	"github.com/rivo/tview"
)

type BuildCommands struct {
	b            *BuildUI
	bc           *BottomCommandsUI
	state        *BuildCommandsState
	runBuildBtn  *tview.Button
	toggleLogBtn *tview.Button
}

type BuildCommandsState struct {
	lbl string
}

func NewBuildCommands(b *BuildUI, bc *BottomCommandsUI) *BuildCommands {
	// View variables
	c := &BuildCommands{
		bc: bc,
		b:  b,
		state: &BuildCommandsState{
			lbl: "",
		},
	}

	// State variables
	var lbl string = "Maximize log panel"

	// Event listeners
	c.toggleLogBtn = bc.AddButton(lbl, 'l', c.onToggleLogClick)
	c.runBuildBtn = bc.AddButton("Build", 's', c.onBuildClick)

	return c
}

// View update functions
func (c *BuildCommands) setToggleLogButton(value string) {
	formatted := c.bc.FormatLabel(value, 'l')
	c.toggleLogBtn.SetLabel(formatted)
}

// State update functions
func (c *BuildCommands) setToggleLogLabel(value string) {
	if c.state.lbl != value {
		c.state.lbl = value
		c.setToggleLogButton(value)
	}
}

func (c *BuildCommands) SetToggleLogLabelBasedOnState(min bool) {
	if min {
		c.setToggleLogLabel("Collapse log panel")
	} else {
		c.setToggleLogLabel("Expand log panel")
	}
}

// Event listeners
func (c *BuildCommands) onToggleLogClick() {
	//ds := c.ds
	//ds.toggleServerRunning()
	//ds.ui.SetFocus(ds.info)
}

func (c *BuildCommands) onBuildClick() {
	c.b.buildActor.StartBuild()
}

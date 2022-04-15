package ui

import "github.com/rivo/tview"

type DevServerCommands struct {
	ds     *DevServerUI
	bc     *BottomCommandsUI
	state  *devServerCommandsState
	svrBtn *tview.Button
}

type devServerCommandsState struct {
	lbl string
}

func NewDevServerCommands(ds *DevServerUI, bc *BottomCommandsUI) *DevServerCommands {
	// View variables
	c := &DevServerCommands{
		bc: bc,
		ds: ds,
		state: &devServerCommandsState{
			lbl: "",
		},
	}

	// State variables
	var lbl string = "Start server"

	// Event listeners
	c.svrBtn = bc.AddButton(lbl, 's', c.onServerClick)

	return c
}

// View update functions
func (c *DevServerCommands) setServerButtonLabel(value string) {
	formatted := c.bc.FormatLabel(value, 's')
	c.svrBtn.SetLabel(formatted)
}

// State update functions
func (c *DevServerCommands) setServerLabel(value string) {
	if c.state.lbl != value {
		c.state.lbl = value
		c.setServerButtonLabel(value)
	}
}

func (c *DevServerCommands) SetServerLabelBasedOnRunning(running bool) {
	if running {
		c.setServerLabel("Stop server")
	} else {
		c.setServerLabel("Start server")
	}
}

// Event listeners
func (c *DevServerCommands) onServerClick() {
	ds := c.ds
	ds.toggleServerRunning()
	ds.ui.SetFocus(ds.info)
}

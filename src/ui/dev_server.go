package ui

import (
	"os/exec"
	"syscall"

	"astroterm/src/astro"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DevServerUI struct {
	*tview.Flex
	app   *tview.Application
	tv    *tview.TextView
	state *serverState
}

type serverState struct {
	running bool
	cmd     *exec.Cmd
}

func NewDevServer(app *tview.Application) *DevServerUI {
	flex := tview.NewFlex()
	flex.SetTitle("Development Server")
	flex.SetTitleAlign(tview.AlignLeft)
	flex.SetBorder(true)
	flex.SetDirection(tview.FlexRow)

	state := &serverState{
		running: false,
		cmd:     nil,
	}

	tv := tview.NewTextView()
	tv.Write([]byte("This is some initial text\n"))
	tv.SetChangedFunc(func() {
		app.Draw()
	})

	devServer := &DevServerUI{
		Flex:  flex,
		app:   app,
		tv:    tv,
		state: state,
	}

	var btn *tview.Button
	form := tview.NewForm()
	form.AddButton("Start", nil)
	form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	btn = form.GetButton(0)
	btn.SetSelectedFunc(func() {
		var label string
		if state.running {
			label = "Start"
		} else {
			label = "Stop"
		}
		btn.SetLabel(label)
		state.running = !state.running
		app.SetFocus(flex)

		if state.running {
			devServer.startServer()
			form.SetButtonBackgroundColor(tcell.ColorDarkRed)
		} else {
			devServer.killServer()
		}
	})
	btn.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDown:
			if btn.InRect(event.Position()) {
				form.SetButtonBackgroundColor(Styles.MoreContrastBackgroundColor)
				go (func() {
					app.Draw()
				})()
			}
		case tview.MouseLeftUp:
			if btn.InRect(event.Position()) {
				form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
			}
		}

		return action, event
	})
	flex.AddItem(form, 3, 0, false)
	flex.AddItem(tv, 0, 1, false)

	return devServer
}

func (ds *DevServerUI) startServer() error {
	cmd, err := astro.RunCommand(astro.Dev, ds.tv)
	if err != nil {
		return err
	}
	ds.state.cmd = cmd
	return nil
}

func (ds *DevServerUI) killServer() error {
	state := ds.state
	if state.cmd != nil {
		err := state.cmd.Process.Signal(syscall.SIGKILL)
		state.cmd.Process.Wait()
		return err
	}
	return nil
}

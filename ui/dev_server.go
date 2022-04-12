package ui

import (
	"os"
	"syscall"

	"astroterm/astro"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DevServerUI struct {
	Flex  *tview.Flex
	app   *tview.Application
	tv    *tview.TextView
	state *serverState
}

type serverState struct {
	running bool
	proc    *os.Process
	pid     int
}

func NewDevServer(app *tview.Application) *DevServerUI {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	state := &serverState{
		running: false,
		proc:    nil,
	}

	tv := tview.NewTextView()
	tv.SetTitle("Logs")
	tv.SetTitleAlign(tview.AlignLeft)
	tv.SetBorder(true)
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
	form.SetTitle("Info")
	form.SetTitleAlign(tview.AlignLeft)
	form.SetBorder(true)
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
	MakeToggleableButton(btn, form, app)
	flex.AddItem(form, 5, 0, false)
	flex.AddItem(tv, 0, 1, false)

	return devServer
}

func (ds *DevServerUI) Primitive() tview.Primitive {
	return ds.Flex
}

func (ds *DevServerUI) Stop() bool {
	ds.killServer()
	return true
}

func (ds *DevServerUI) startServer() error {
	cmd, err := astro.RunCommand(astro.Dev, ds.tv)
	if err != nil {
		return err
	}
	ds.state.proc = cmd.Process
	ds.state.pid = cmd.Process.Pid
	return nil
}

func (ds *DevServerUI) killServer() error {
	state := ds.state
	if state.proc != nil {
		childPid := state.pid + 1
		childProc, childErr := os.FindProcess(childPid)
		if childErr == nil {
			childProc.Signal(syscall.SIGKILL)
			childProc.Wait()
		}

		err := state.proc.Signal(syscall.SIGKILL)
		state.proc.Wait()

		// If there was a child then this might error
		if err != nil && childErr == nil {
			return err
		}

		return nil
	}
	return nil
}

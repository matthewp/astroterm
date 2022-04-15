package ui

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"syscall"

	"astroterm/astro"
	"astroterm/db"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DevServerUI struct {
	Flex      *tview.Flex
	ui        *UI
	info      *tview.Flex
	logs      *tview.TextView
	ovw       *tview.TextView
	state     *serverState
	focusMenu func()
}

type serverState struct {
	running bool
}

var portMatch = regexp.MustCompile("(localhost|127.0.0.1):([0-9]{4})\\/")

func NewDevServer(u *UI) *DevServerUI {
	var devServer *DevServerUI
	var info *tview.Flex
	// Start server button
	var btn *tview.Button

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.SetFocusFunc(func() {
		u.SetFocus(info)
	})

	state := &serverState{
		running: false,
	}

	logs := tview.NewTextView()
	logs.SetTitle("Logs")
	logs.SetTitleAlign(tview.AlignLeft)
	logs.SetBorder(true)
	logs.SetChangedFunc(func() {
		u.Draw()
	})
	logs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyLeft:
			devServer.focusMenu()
			return nil
		case tcell.KeyUp:
		case tcell.KeyBacktab:
			u.SetFocus(info)
			return nil
		}
		return event
	})

	// The Info Section
	info = tview.NewFlex()
	info.SetTitle("Info")
	info.SetTitleAlign(tview.AlignLeft)
	info.SetBorder(true)
	info.SetBorderPadding(0, 0, 1, 1)
	info.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyLeft:
		case tcell.KeyBacktab:
			devServer.focusMenu()
			return nil
		case tcell.KeyDown:
			u.SetFocus(logs)
			return nil
		case tcell.KeyTab:
			if u.GetFocus() == btn {
				u.SetFocus(logs)
			} else {
				u.SetFocus(btn)
			}
			return nil
		}

		return event

	})

	// Overview info
	ovwf := tview.NewFlex()
	ovwf.SetDirection(tview.FlexRow)

	ovw := tview.NewTextView()
	ovwf.AddItem(nil, 0, 1, false).
		AddItem(ovw, 0, 1, false).
		AddItem(nil, 0, 1, false)

	devServer = &DevServerUI{
		Flex:  flex,
		ui:    u,
		info:  info,
		logs:  logs,
		ovw:   ovw,
		state: state,
	}
	devServer.LoadDeverServerModel()
	if devServer.Model().Pid != 0 {
		state.running = true
	}

	form := tview.NewForm()
	form.AddButton("Start server", nil)
	form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	form.SetButtonsAlign(tview.AlignCenter)
	btn = form.GetButton(0)
	SetButtonState(state, btn, form)

	btn.SetSelectedFunc(func() {
		SetButtonState(state, btn, form)
		u.SetFocus(flex)
		devServer.toggleServerState()
	})
	MakeToggleableButton(btn, form, u)

	info.AddItem(ovwf, 0, 1, false)
	info.AddItem(form, 0, 1, false)

	flex.AddItem(info, 5, 0, false)
	flex.AddItem(logs, 0, 1, false)

	// Set initial state
	devServer.setOverviewInformation()

	return devServer
}

func (ds *DevServerUI) Primitive() tview.Primitive {
	return ds.Flex
}

func (ds *DevServerUI) Stop() {
	ds.shutdownServer(false)
}

func (ds *DevServerUI) SetFocusMenu(focusMenu func()) {
	ds.focusMenu = focusMenu
}

func (ds *DevServerUI) MakeActive(cmds *BottomCommandsUI) {
	var lbl string
	if ds.state.running {
		lbl = "Stop server"
	} else {
		lbl = "Start server"
	}
	cmds.AddButton(lbl, 's', func() {
		ds.toggleServerState()
		ds.ui.SetFocus(ds.info)
	})
}

func (ds *DevServerUI) Write(p []byte) (int, error) {
	if ds.Model().Port == 0 {
		ds.parseHostInformation(p)
	}
	return ds.logs.Write(p)
}

func (ds *DevServerUI) parseHostInformation(p []byte) {
	part := string(p)
	rs := portMatch.FindStringSubmatch(part)
	if len(rs) > 1 {
		portString := rs[2]
		ds.Model().Port, _ = strconv.Atoi(portString)
		ds.Model().Hostname = rs[1]
		ds.setOverviewInformation()
		err := ds.ui.db.SetDevServerInformation(ds.Model())
		if err != nil {
			ds.logs.Write([]byte(err.Error()))
		}
	}
}

func (ds *DevServerUI) setOverviewInformation() {
	model := ds.Model()

	var msg string
	if model.Hostname == "" {
		msg = "No server running"
	} else {
		msg = fmt.Sprintf("Listening at http://%s:%v", model.Hostname, model.Port)
	}
	ds.ovw.SetText(msg)
}

func (ds *DevServerUI) toggleServerState() {
	state := ds.state
	state.running = !state.running
	if state.running {
		ds.startServer()
	} else {
		ds.shutdownServer(true)
	}
}

func (ds *DevServerUI) startServer() error {
	cmd, err := astro.RunCommand(astro.Dev, ds)
	if err != nil {
		return err
	}
	ds.Model().Pid = cmd.Process.Pid
	ds.ui.db.AddStartedDevServer(ds.Model()) // TODO change
	return nil
}

func (ds *DevServerUI) shutdownServer(updateUI bool) error {
	e1 := ds.killServer()
	ds.Model().Hostname = ""
	ds.Model().Port = 0
	if updateUI {
		ds.setOverviewInformation()
	}
	e2 := ds.ui.db.DeleteDevServer(ds.Model())
	// This must happen after the model is deleted from the database
	ds.Model().Pid = 0

	if e1 != nil {
		return e1
	}
	return e2
}

func (ds *DevServerUI) killServer() error {
	if ds.Model().Pid != 0 {
		err1 := KillPid(ds.Model().Pid + 1)
		err2 := KillPid(ds.Model().Pid)

		if err1 != nil {
			return err1
		}
		return err2
	}
	return nil
}

func KillPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err == nil {
		proc.Signal(syscall.SIGKILL)
		proc.Wait()
	}
	return err
}

func (ds *DevServerUI) Model() *db.DevServerModel {
	return ds.ui.DevModel
}

func (ds *DevServerUI) LoadDeverServerModel() {
	projectDir := ds.ui.CurrentProject.Dir
	model, err := ds.ui.db.LoadDeverServerModel(projectDir)
	if err != nil {
		ds.ui.DevModel.ProjectDir = projectDir
		return
	}
	if model != nil {
		ds.ui.DevModel = model
	} else {
		ds.ui.DevModel.ProjectDir = projectDir
	}
}

func SetButtonState(state *serverState, btn *tview.Button, form *tview.Form) {
	if state.running {
		btn.SetLabel("Stop server")
		form.SetButtonBackgroundColor(tcell.ColorDarkRed)
	} else {
		btn.SetLabel("Start server")
		form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	}

}

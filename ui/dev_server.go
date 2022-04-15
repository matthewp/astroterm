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
	Flex  *tview.Flex
	ui    *UI
	form  *tview.Form
	info  *tview.Flex
	logs  *tview.TextView
	ovw   *tview.TextView
	ssBtn *tview.Button
	cmds  *DevServerCommands

	state     *serverState
	focusMenu func()
}

type serverState struct {
	running           bool
	active            bool
	serverButtonColor tcell.Color
	serverButtonText  string
}

var portMatch = regexp.MustCompile("(localhost|127.0.0.1):([0-9]{4})\\/")

func NewDevServer(u *UI) *DevServerUI {
	var ds *DevServerUI

	// Views
	var info *tview.Flex
	var ssBtn *tview.Button
	form := tview.NewForm()
	flex := tview.NewFlex()
	info = tview.NewFlex()
	logs := tview.NewTextView()
	ovwf := tview.NewFlex()
	ovw := tview.NewTextView()

	// State
	LoadDevServerModel(u)
	state := &serverState{
		running:           u.DevModel.IsRunning(),
		active:            false,
		serverButtonText:  "",
		serverButtonColor: tcell.ColorBlack,
	}
	ds = &DevServerUI{
		Flex:  flex,
		ui:    u,
		form:  form,
		info:  info,
		logs:  logs,
		ovw:   ovw,
		ssBtn: ssBtn,
		state: state,
	}

	// Event listeners
	flex.SetFocusFunc(func() {
		u.SetFocus(info)
	})
	logs.SetChangedFunc(func() {
		u.Draw()
	})
	logs.SetInputCapture(ds.logInputCapture)
	info.SetInputCapture(ds.infoInputCapture)
	form.AddButton("Start server", nil)
	ssBtn = form.GetButton(0)
	ds.ssBtn = ssBtn
	ssBtn.SetSelectedFunc(ds.ssBtnSelected)
	MakeToggleableButton(ssBtn, form, u)

	// Initialization

	form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	form.SetButtonsAlign(tview.AlignCenter)

	flex.SetDirection(tview.FlexRow)
	flex.AddItem(info, 5, 0, false)
	flex.AddItem(logs, 0, 1, false)

	logs.SetTitle("Logs")
	logs.SetTitleAlign(tview.AlignLeft)
	logs.SetBorder(true)

	info.SetTitle("Info")
	info.SetTitleAlign(tview.AlignLeft)
	info.SetBorder(true)
	info.SetBorderPadding(0, 0, 1, 1)
	info.AddItem(ovwf, 0, 1, false)
	info.AddItem(form, 0, 1, false)

	ovwf.SetDirection(tview.FlexRow)
	ovwf.AddItem(nil, 0, 1, false).
		AddItem(ovw, 0, 1, false).
		AddItem(nil, 0, 1, false)

	return ds
}

// Implementations
func (ds *DevServerUI) Primitive() tview.Primitive {
	return ds.Flex
}

func (ds *DevServerUI) Stop() {
	// TODO UI?
	ds.setActive(false)
	ds.setServerRunning(false)
}

func (ds *DevServerUI) SetFocusMenu(focusMenu func()) {
	ds.focusMenu = focusMenu
}

func (ds *DevServerUI) MakeActive(cmds *BottomCommandsUI) {
	ds.setActive(true)

	setServerButtonColorBasedOnRunning(ds)
	setServerButtonTextBasedOnRunning(ds)
	setOverviewText(ds)

	ds.cmds = NewDevServerCommands(ds, cmds)
	ds.cmds.SetServerLabelBasedOnRunning(ds.state.running)
}

func (ds *DevServerUI) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 's':
		ds.toggleServerRunning()
		break
	}
	return event
}

func (ds *DevServerUI) Write(p []byte) (int, error) {
	if ds.Model().Port == 0 {
		hostname, port := ds.parseHostInformation(p)
		if hostname != "" {
			ds.setHostAndPort(hostname, port)
			saveDevServerInformation(ds)
		}
	}
	return appendLogText(ds, p)
}

// View update functions
func setOverviewText(ds *DevServerUI) {
	model := ds.Model()
	var msg string
	if model.Hostname == "" {
		msg = "No server running"
	} else {
		msg = fmt.Sprintf("Listening at http://%s:%v", model.Hostname, model.Port)
	}
	ds.ovw.SetText(msg)
}

func appendLogText(ds *DevServerUI, bytes []byte) (int, error) {
	return ds.logs.Write(bytes)
}

func setServerButtonLabel(ds *DevServerUI, value string) {
	ds.ssBtn.SetLabel(value)
}

func setFormButtonBackgroundColor(ds *DevServerUI, value tcell.Color) {
	ds.form.SetButtonBackgroundColor(value)
}

// State update functions
func setServerButtonText(ds *DevServerUI, value string) {
	if ds.state.serverButtonText != value {
		ds.state.serverButtonText = value
		setServerButtonLabel(ds, value)
	}
}

func setServerButtonTextBasedOnRunning(ds *DevServerUI) {
	if ds.state.running {
		setServerButtonText(ds, "Stop server")
	} else {
		setServerButtonText(ds, "Start server")
	}
}

func setServerButtonColor(ds *DevServerUI, value tcell.Color) {
	if ds.state.serverButtonColor != value {
		ds.state.serverButtonColor = value
		setFormButtonBackgroundColor(ds, value)
	}
}

func setServerButtonColorBasedOnRunning(ds *DevServerUI) {
	if ds.state.running {
		setServerButtonColor(ds, tcell.ColorDarkRed)
	} else {
		setServerButtonColor(ds, Styles.ContrastBackgroundColor)
	}
}

func (ds *DevServerUI) setHostAndPort(hostname string, port int) {
	ds.Model().Port = port
	ds.Model().Hostname = hostname
	setOverviewText(ds)
}

func (ds *DevServerUI) setPid(value int) {
	ds.Model().Pid = value
}

func (ds *DevServerUI) setServerRunning(value bool) {
	state := ds.state
	if state.running != value {
		state.running = value
		if state.running {
			if state.active {
				setServerButtonTextBasedOnRunning(ds)
				setServerButtonColorBasedOnRunning(ds)
				ds.cmds.SetServerLabelBasedOnRunning(true)
			}

			ds.startServer()
		} else {
			if state.active {
				setServerButtonTextBasedOnRunning(ds)
				setServerButtonColorBasedOnRunning(ds)
				ds.setHostAndPort("", 0)
				ds.cmds.SetServerLabelBasedOnRunning(false)
			}

			ds.shutdownServer()

			// This must happen after the model is deleted from the database
			ds.setPid(0)
		}
	}
}

func (ds *DevServerUI) toggleServerRunning() {
	ds.setServerRunning(!ds.state.running)
}

func (ds *DevServerUI) setActive(value bool) {
	ds.state.active = value
}

// Event listeners
func (ds *DevServerUI) logInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
	case tcell.KeyLeft:
		ds.focusMenu()
		return nil
	case tcell.KeyUp, tcell.KeyBacktab:
		ds.ui.SetFocus(ds.info)
		return nil
	}
	return event
}

func (ds *DevServerUI) infoInputCapture(event *tcell.EventKey) *tcell.EventKey {
	u := ds.ui
	switch key := event.Key(); key {
	case tcell.KeyLeft, tcell.KeyBacktab:
		ds.focusMenu()
		return nil
	case tcell.KeyDown:
		u.SetFocus(ds.logs)
		return nil
	case tcell.KeyTab:
		if u.GetFocus() == ds.ssBtn {
			u.SetFocus(ds.logs)
		} else {
			u.SetFocus(ds.ssBtn)
		}
		return nil
	}

	return event
}

func (ds *DevServerUI) ssBtnSelected() {
	u := ds.ui
	u.SetFocus(ds.Flex)
	ds.toggleServerRunning()
}

// Logic functions
func (ds *DevServerUI) parseHostInformation(p []byte) (string, int) {
	part := string(p)
	rs := portMatch.FindStringSubmatch(part)
	if len(rs) > 1 {
		portString := rs[2]
		port, _ := strconv.Atoi(portString)
		hostname := rs[1]
		return hostname, port
	}
	return "", 0
}

func (ds *DevServerUI) Model() *db.DevServerModel {
	return ds.ui.DevModel
}

// Effects
func saveDevServerInformation(ds *DevServerUI) {
	err := ds.ui.DB.SetDevServerInformation(ds.Model())
	if err != nil {
		appendLogText(ds, []byte(err.Error()))
	}
}

func (ds *DevServerUI) startServer() error {
	cmd, err := astro.RunCommand(astro.Dev, ds)
	if err != nil {
		return err
	}
	ds.setPid(cmd.Process.Pid)
	ds.ui.DB.AddStartedDevServer(ds.Model())
	return nil
}

func (ds *DevServerUI) killServer() error {
	if ds.Model().Pid != 0 {
		err1 := killPid(ds.Model().Pid + 1)
		err2 := killPid(ds.Model().Pid)

		if err1 != nil {
			return err1
		}
		return err2
	}
	return nil
}

func (ds *DevServerUI) shutdownServer() error {
	e1 := ds.killServer()
	e2 := ds.ui.DB.DeleteDevServer(ds.Model())
	if e1 != nil {
		return e1
	}
	return e2
}

func killPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err == nil {
		proc.Signal(syscall.SIGKILL)
		proc.Wait()
	}
	return err
}

// TODO refactor

func LoadDevServerModel(u *UI) {
	projectDir := u.CurrentProject.Dir
	model, err := u.DB.LoadDevServerModel(projectDir)
	if err != nil {
		u.DevModel.ProjectDir = projectDir
		return
	}
	if model != nil {
		u.DevModel = model
	} else {
		u.DevModel.ProjectDir = projectDir
	}
}

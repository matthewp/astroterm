package ui

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"syscall"

	"astroterm/astro"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DevServerUI struct {
	Flex  *tview.Flex
	app   *tview.Application
	logs  *tview.TextView
	ovw   *tview.TextView
	state *serverState
}

type serverState struct {
	running  bool
	proc     *os.Process
	pid      int
	hostname string
	port     int
}

var portMatch = regexp.MustCompile("(localhost|127.0.0.1):([0-9]{4})\\/")

func NewDevServer(app *tview.Application) *DevServerUI {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	state := &serverState{
		running:  false,
		proc:     nil,
		pid:      0,
		port:     0,
		hostname: "",
	}

	logs := tview.NewTextView()
	logs.SetTitle("Logs")
	logs.SetTitleAlign(tview.AlignLeft)
	logs.SetBorder(true)
	logs.SetChangedFunc(func() {
		app.Draw()
	})

	// The Info Section
	info := tview.NewFlex()
	info.SetTitle("Info")
	info.SetTitleAlign(tview.AlignLeft)
	info.SetBorder(true)
	info.SetBorderPadding(0, 0, 1, 1)

	// Overview info
	ovwf := tview.NewFlex()
	ovwf.SetDirection(tview.FlexRow)

	ovw := tview.NewTextView()
	ovwf.AddItem(nil, 0, 1, false).
		AddItem(ovw, 0, 1, false).
		AddItem(nil, 0, 1, false)

	devServer := &DevServerUI{
		Flex:  flex,
		app:   app,
		logs:  logs,
		ovw:   ovw,
		state: state,
	}

	var btn *tview.Button
	form := tview.NewForm()
	form.AddButton("Start", nil)
	form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
	form.SetButtonsAlign(tview.AlignCenter)
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
			state.hostname = ""
			state.port = 0
			devServer.setOverviewInformation()
		}
	})
	MakeToggleableButton(btn, form, app)

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
	ds.killServer()
}

func (ds *DevServerUI) Write(p []byte) (int, error) {
	if ds.state.port == 0 {
		ds.parseHostInformation(p)
	}
	return ds.logs.Write(p)
}

func (ds *DevServerUI) parseHostInformation(p []byte) {
	part := string(p)
	rs := portMatch.FindStringSubmatch(part)
	if len(rs) > 1 {
		portString := rs[2]
		ds.state.port, _ = strconv.Atoi(portString)
		ds.state.hostname = rs[1]
		ds.setOverviewInformation()
	}
}

func (ds *DevServerUI) setOverviewInformation() {
	state := ds.state

	var msg string
	if state.hostname == "" {
		msg = "No server running"
	} else {
		msg = fmt.Sprintf("Listening at http://%s:%v", state.hostname, state.port)
	}
	ds.ovw.SetText(msg)
}

func (ds *DevServerUI) startServer() error {
	cmd, err := astro.RunCommand(astro.Dev, ds)
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

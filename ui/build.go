package ui

import (
	"astroterm/actors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BuildUI struct {
	Flex *tview.Flex

	buildActor *actors.BuildActor
	cmds       *BuildCommands
	focusMenu  func()
	logs       *tview.TextView
	overview   *tview.Flex
	ui         *UI
}

func NewBuildUI(u *UI, buildActor *actors.BuildActor) *BuildUI {
	b := &BuildUI{
		buildActor: buildActor,
		ui:         u,
	}

	flex := tview.NewFlex()
	b.Flex = flex
	flex.SetDirection(tview.FlexRow)

	overview := tview.NewFlex()
	b.overview = overview
	overview.SetBorder(true)
	overview.SetTitle("Overview")
	overview.SetTitleAlign(tview.AlignLeft)
	flex.AddItem(overview, 0, 1, false)

	logs := tview.NewTextView()
	b.logs = logs
	flex.AddItem(logs, 2, 0, false)

	// Event listeners
	flex.SetFocusFunc(func() {
		u.SetFocus(overview)
	})
	overview.SetInputCapture(b.overviewInputCapture)
	logs.SetInputCapture(b.logInputCapture)
	logs.SetFocusFunc(b.logFocusFunc)

	// Initialization
	logs.SetTitle("Logs")
	logs.SetTitleAlign(tview.AlignLeft)
	logs.SetBorder(true)

	go b.listenForBuildEvents()

	return b
}

func (b *BuildUI) listenForBuildEvents() {
	lchan := b.buildActor.SubscribeToLogs()

	for {
		data := <-lchan
		b.appendLogText(data)
	}
}

func (b *BuildUI) appendLogText(content string) (int, error) {
	bytes := []byte(content)
	return b.logs.Write(bytes)
}

/* UISection implementation */
func (b *BuildUI) Primitive() tview.Primitive {
	return b.Flex
}

func (b *BuildUI) Stop() {}

func (b *BuildUI) SetFocusMenu(focusMenu func()) {
	b.focusMenu = focusMenu
}

func (b *BuildUI) MakeActive(cmds *BottomCommandsUI) {
	b.cmds = NewBuildCommands(b, cmds)
}

func (b *BuildUI) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 's':
		b.cmds.onBuildClick()
		break
	}
	return event
}

// Logic functions
func (b *BuildUI) expandLogs() {
	b.Flex.ResizeItem(b.logs, 0, 1)
}

func (b *BuildUI) collapseLogs() {
	b.Flex.ResizeItem(b.logs, 2, 0)
}

// Event listeners
func (b *BuildUI) overviewInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
	case tcell.KeyLeft, tcell.KeyBacktab:
		b.focusMenu()
		return nil
	case tcell.KeyDown, tcell.KeyTab:
		b.ui.SetFocus(b.logs)
		return nil
	}
	return event
}

func (b *BuildUI) logInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
	case tcell.KeyLeft:
		b.collapseLogs()
		b.focusMenu()
		return nil
	case tcell.KeyUp, tcell.KeyBacktab:
		b.collapseLogs()
		b.ui.SetFocus(b.overview)
		return nil
	}

	return event
}

func (b *BuildUI) logFocusFunc() {
	b.expandLogs()
}

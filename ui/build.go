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
	logs       *tview.TextView
}

func NewBuildUI(u *UI, buildActor *actors.BuildActor) *BuildUI {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	overview := tview.NewFlex()
	overview.SetBorder(true)
	overview.SetTitle("Overview")
	flex.AddItem(overview, 0, 1, false)

	logs := tview.NewTextView()
	flex.AddItem(logs, 0, 1, false)

	// Initialization
	logs.SetTitle("Logs")
	logs.SetTitleAlign(tview.AlignLeft)
	logs.SetBorder(true)

	b := &BuildUI{
		Flex: flex,

		buildActor: buildActor,
		logs:       logs,
	}

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

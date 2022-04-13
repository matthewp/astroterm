package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MakeToggleableButton(btn *tview.Button, form *tview.Form, u *UI) {
	btn.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDown:
			if btn.InRect(event.Position()) {
				form.SetButtonBackgroundColor(Styles.MoreContrastBackgroundColor)
				go (func() {
					u.Draw()
				})()
			}
		case tview.MouseLeftUp:
			if btn.InRect(event.Position()) {
				form.SetButtonBackgroundColor(Styles.ContrastBackgroundColor)
			}
		}

		return action, event
	})
}

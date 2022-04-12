package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
	labelColor:            Styles.SecondaryTextColor,
	fieldBackgroundColor:  Styles.ContrastBackgroundColor,
	fieldTextColor:        Styles.PrimaryTextColor,
	buttonBackgroundColor: Styles.ContrastBackgroundColor,
	buttonTextColor:       Styles.PrimaryTextColor,

			button.SetLabelColor(f.buttonTextColor).
		SetLabelColorActivated(f.buttonBackgroundColor).
		SetBackgroundColorActivated(f.buttonTextColor).
		SetBackgroundColor(f.buttonBackgroundColor)


						form.SetButtonBackgroundColor(tcell.ColorOrchid)
				form.SetButtonTextColor(tcell.ColorPeachPuff)
*/

var Styles = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorSeaGreen,
	MoreContrastBackgroundColor: tcell.ColorOrchid,
	BorderColor:                 tcell.ColorWhite,
	TitleColor:                  tcell.ColorWhite,
	GraphicsColor:               tcell.ColorWhite,
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
}

var NavStyles = struct {
	BackgroundColor tcell.Color
	TextColor       tcell.Color
}{
	tcell.ColorLightGray,
	tcell.ColorBlack,
}

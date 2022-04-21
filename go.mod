module astroterm

go 1.18

replace "github.com/matthewp/bestbar" => "../bestbar"

require (
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/matthewp/bestbar v0.2.0
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/rivo/tview v0.0.0-20220307222120-9994674d60a8
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.0.0-20220318055525-2edf467146b5 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.7 // indirect
)

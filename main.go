package main

import (
	"astroterm/cli"
	"astroterm/ui"
)

func main() {
	if cli.Init() {
		startUI()
	}
}

func startUI() {
	u := ui.NewUI()
	if err := u.Start(); err != nil {
		panic(err)
	}
}

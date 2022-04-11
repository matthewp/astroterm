package main

import (
	"astroterm/src/env"
	"astroterm/src/ui"
	"fmt"
)

func main() {
	start()
}

func testing() {
	env, err := env.GetEnvironment()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PWD %s", env.ConfigPath)
}

func start() {
	u := ui.NewUI()
	if err := u.Start(); err != nil {
		panic(err)
	}
}

package cli

import (
	"astroterm/actors"
	"astroterm/project"
	"fmt"
	"log"
)

// astroterm build
func RunBuild() {
	proj, err := project.OpenLocalProject()
	if err != nil {
		log.Fatal(err)
		return
	}
	buildActor := actors.NewBuildActor(proj)
	err = <-buildActor.RunBuildToCompletion()
	fmt.Printf("Got an error? %v\n", err)
	if err != nil {
		log.Fatal(err)
	}
	return
}

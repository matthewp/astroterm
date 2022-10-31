package cli

import (
	"astroterm/home"
	"astroterm/project"
	"errors"
	"fmt"
	"os"
)

func Info(args []string) error {
	if len(args) == 0 {
		fmt.Printf("astroterm info expects a subcommand")
		return nil
	}

	n := args[0]
	switch n {
	case "dev-log-file":
		proj, err := project.OpenLocalProject()
		if err != nil {
			return err
		}

		logpth, err := home.GetDevLogPath(proj.Dir)
		if err != nil {
			return err
		}

		if _, err := os.Stat(logpth); errors.Is(err, os.ErrNotExist) {
			// logpth does not exist
			return nil
		}

		fmt.Println(logpth)

		return nil
	default:
		fmt.Printf("Unknown subcommand %v\n", n)
		return nil
	}
}

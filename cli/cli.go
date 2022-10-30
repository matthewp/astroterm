package cli

import (
	"astroterm/version"
	"flag"
	"fmt"
	"os"
)

///https://github.com/benweidig/tortuga/blob/master/Makefile

func Init() bool {
	arg1 := ""
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	}

	switch arg1 {
	case "run-script":
		rsCmd := flag.NewFlagSet("run-script", flag.ExitOnError)
		var pipe string
		rsCmd.StringVar(&pipe, "pipe", "", "Pipe output to a file")
		var name string
		rsCmd.StringVar(&name, "name", "", "Name of the script to run")
		rsCmd.Parse(os.Args[2:])
		RunScript(name, pipe)
		return false

	default:
		var version bool
		versionUsage := "Show the version and exit"
		flag.BoolVar(&version, "version", false, versionUsage)
		flag.BoolVar(&version, "v", false, shorthandUsage(versionUsage))
		var help bool
		helpUsage := "Show the help message and exit"
		flag.BoolVar(&help, "help", false, helpUsage)
		flag.BoolVar(&help, "h", false, shorthandUsage((helpUsage)))

		flag.Parse()

		if help {
			Usage()
			return false
		}

		if version {
			Version()
			return false
		}
		break
	}

	return true
}

func shorthandUsage(usage string) string {
	return usage + " (shorthand)"
}

func Usage() {
	v := version.Version
	fmt.Printf(`astroterm %s

USAGE:
	astroterm [FLAGS] [OPTIONS]

FLAGS:
	-h, --help		Prints help information
	-v, --version		Prints version information
`, v)
}

func Version() {
	v := version.Version
	fmt.Printf("%v\n", v)
}

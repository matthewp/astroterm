package astro

import (
	"os"
	"os/exec"
)

func runScript(scriptName string, logpath string) (int, string, error) {
	astrotermBin := os.Args[0]
	cmd := exec.Command(astrotermBin, "run-script", "--name", scriptName, "--pipe", logpath)
	err := cmd.Start()
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, logpath, nil
}

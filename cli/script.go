package cli

import (
	"astroterm/util"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Run an npm script
func RunScript(scriptName string, pipeTo string) error {
	cmd := exec.Command("npm", "run", scriptName)

	var err error
	if pipeTo != "" {
		// open the out file for writing
		outfile, err := os.Create(pipeTo)
		if err != nil {
			return err
		}
		defer outfile.Close()
		cmd.Stdout = outfile
		cmd.Stderr = outfile
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	pid := cmd.Process.Pid

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		util.TermPid(pid)
		util.TermPid(pid + 1) // Kill the second npm run dev
		os.Exit(0)
	}()

	cmd.Wait()
	return nil
}

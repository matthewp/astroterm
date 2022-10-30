package cli

import (
	"os"
	"os/exec"
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
	}

	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

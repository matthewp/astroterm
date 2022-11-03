package astro

import (
	"astroterm/home"
	"os/exec"
)

func RunBuildAndPipeToLog(projectDir string) (int, string, error) {
	logpath, err := home.GetBuildLogPath(projectDir)
	if err != nil {
		return 0, "", err
	}
	return runScriptAndPipeToLog("build", logpath)
}

func RunBuild(projectDir string) (*exec.Cmd, error) {
	cmd, err := runScript("build")
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

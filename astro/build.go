package astro

import (
	"astroterm/home"
)

func RunBuildAndPipeToLog(projectDir string) (int, string, error) {
	logpath, err := home.GetBuildLogPath(projectDir)
	if err != nil {
		return 0, "", err
	}
	return runScript("build", logpath)
}

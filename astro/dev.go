package astro

import (
	"astroterm/home"
)

func RunDevAndPipeToLog(projectDir string) (int, string, error) {
	logpath, err := home.GetDevLogPath(projectDir)
	if err != nil {
		return 0, "", err
	}
	return runScriptAndPipeToLog("dev", logpath)
}

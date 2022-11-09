package astro

import (
	aenv "astroterm/env"
	"astroterm/home"
	"fmt"
	"os"
	"os/exec"
)

func RunBuildAndPipeToLog(projectDir string) (int, string, error) {
	logpath, err := home.GetBuildLogPath(projectDir)
	if err != nil {
		return 0, "", err
	}
	return runScriptAndPipeToLog("build", logpath)
}

func CreateBuildCommandWithCustomConfig(projectDir string, configSource []byte) (*exec.Cmd, error) {
	env, err := aenv.GetEnvironment()
	if err != nil {
		return nil, err
	}

	binPath, err := getAstroBinaryPath(env)
	if err != nil {
		return nil, err
	}

	f, err := os.CreateTemp("", "astro.*.config.mjs")
	fmt.Printf("Got temp file %v\n", f.Name())
	if err != nil {
		return nil, err
	}

	_, err = f.Write(configSource)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("node", binPath, "build", "--root", projectDir, "--config", f.Name())

	return cmd, nil
}

func RunBuild(projectDir string) (*exec.Cmd, error) {
	cmd, err := runScript("build")
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

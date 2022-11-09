package astro

import (
	aenv "astroterm/env"
	"errors"
	"io"
	"os/exec"
	"path"
)

type AstroCommand string

const (
	Dev     AstroCommand = "dev"
	Build   AstroCommand = "build"
	Preview AstroCommand = "preview"
)

func pipetotext(pipe io.ReadCloser, wrtr io.Writer) {
	for {
		if _, err := io.Copy(wrtr, pipe); err != nil {
			break
		}
	}
}

func getAstroBinaryPath(env *aenv.TermEnvironment) (string, error) {
	relBinPath := path.Join("node_modules", ".bin", "astro")
	if aenv.TryFindFile(env.Pwd, relBinPath) == "" {
		return "", errors.New("unable to find the astro binary. Do you need to run npm install?")
	}
	binPath := path.Join(env.Pwd, relBinPath)
	return binPath, nil
}

func RunCommandWithWriter(subcmd AstroCommand, wrtr io.Writer) (*exec.Cmd, error) {
	env, err := aenv.GetEnvironment()
	if err != nil {
		return nil, err
	}

	binPath, err := getAstroBinaryPath(env)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("node", binPath, string(subcmd))

	if wrtr != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}

		go pipetotext(stdout, wrtr)
		go pipetotext(stderr, wrtr)
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

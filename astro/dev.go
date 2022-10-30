package astro

import (
	"astroterm/home"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func GetDevLogPath(projectDir string) (string, error) {
	logdir, err := home.OpenLogDir()
	if err != nil {
		return "", err
	}

	_, dirname := filepath.Split(projectDir)

	hasher := sha1.New()
	hasher.Write([]byte(logdir))
	sha := hex.EncodeToString(hasher.Sum(nil))

	projectlname := dirname + "-" + sha[:7]

	logpath := path.Join(logdir, projectlname) + ".dev.log"
	return logpath, nil
}

func RunDevAndPipeToLog(projectDir string) (int, string, error) {
	logpath, err := GetDevLogPath(projectDir)
	if err != nil {
		return 0, "", err
	}

	astrotermBin := os.Args[0]
	cmd := exec.Command(astrotermBin, "run-script", "--name", "dev", "--pipe", logpath)
	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, logpath, nil
}

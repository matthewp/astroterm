package home

import (
	"crypto/sha1"
	"encoding/hex"
	"path"
	"path/filepath"
)

const (
	devLog   string = "dev"
	buildLog string = "build"
)

func getLogFilePath(projectDir string, w string) (string, error) {
	logdir, err := OpenLogDir()
	if err != nil {
		return "", err
	}

	_, dirname := filepath.Split(projectDir)

	hasher := sha1.New()
	hasher.Write([]byte(logdir))
	sha := hex.EncodeToString(hasher.Sum(nil))

	projectlname := dirname + "-" + sha[:7]

	logpath := path.Join(logdir, projectlname) + "." + w + ".log"
	return logpath, nil
}

func GetDevLogPath(projectDir string) (string, error) {
	return getLogFilePath(projectDir, devLog)
}

func GetBuildLogPath(projectDir string) (string, error) {
	return getLogFilePath(projectDir, buildLog)
}

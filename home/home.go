package home

import (
	"os"
	"path"
)

const hiddenFolder = ".astroterm"

func GetHomeFolderLocation() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	loc := path.Join(dirname, hiddenFolder)
	return loc, nil
}

func GetLogFolderLocation() (string, error) {
	homedir, err := GetHomeFolderLocation()
	if err != nil {
		return "", err
	}
	logdir := path.Join(homedir, "logs")
	return logdir, nil
}

func OpenHomeFolder() (string, error) {
	dirname, err := GetHomeFolderLocation()
	if err != nil {
		return "", err
	}
	return dirname, os.MkdirAll(dirname, os.ModePerm)
}

func OpenLogDir() (string, error) {
	logdir, err := GetLogFolderLocation()
	if err != nil {
		return "", err
	}
	return logdir, os.MkdirAll(logdir, os.ModePerm)
}

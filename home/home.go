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

func OpenHomeFolder() (string, error) {
	dirname, err := GetHomeFolderLocation()
	if err != nil {
		return "", err
	}
	return dirname, os.MkdirAll(dirname, os.ModePerm)
}

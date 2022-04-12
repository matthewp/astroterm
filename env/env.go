package env

import (
	"errors"
	"os"
	"path"
)

type TermEnvironment struct {
	Pwd            string
	ConfigPath     string
	IsAstroProject bool
}

func GetEnvironment() (*TermEnvironment, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var configPath string
	exts := []string{"mjs", "ts"}
	for _, ext := range exts {
		configPath = tryFindConfig(pwd, ext)
		if configPath != "" {
			break
		}
	}

	if configPath == "" {
		pkgjson := tryFindPackageJson(pwd)
		if pkgjson == "" {
			return &TermEnvironment{
				Pwd:            pwd,
				ConfigPath:     configPath,
				IsAstroProject: false,
			}, errors.New("Unable to find an Astro config or a package.json. Is this an Astro project?")
		}

		// TODO support projects without an astro.config.mjs
		return &TermEnvironment{
			Pwd:            pwd,
			ConfigPath:     configPath,
			IsAstroProject: true,
		}, errors.New("Unable to find an Astro config file. Is this an Astro project?")
	}

	return &TermEnvironment{
		Pwd:            pwd,
		ConfigPath:     configPath,
		IsAstroProject: true,
	}, nil
}

func tryFindConfig(pwd string, ext string) string {
	return TryFindFile(pwd, "astro.config."+ext)
}

func tryFindPackageJson(pwd string) string {
	return TryFindFile(pwd, "package.json")
}

func TryFindFile(pwd string, inPth string) string {
	pth := path.Join(pwd, inPth)

	if _, err := os.Stat(pth); errors.Is(err, os.ErrNotExist) {
		return ""
	}

	return pth
}

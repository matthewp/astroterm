package info

import (
	_ "embed"
	"encoding/json"
	"os/exec"
)

//go:embed load-config.mjs
var openConfigSrc string

type IntegrationConfig struct {
	Name string `json:"name"`
}

type BuildConfig struct {
	Format      string `json:"format"`
	Client      string `json:"client"`
	Server      string `json:"server"`
	ServerEntry string `json:"serverEntry"`
}

type ConfigInfo struct {
	SrcDir       string              `json:"srcDir"`
	OutDir       string              `json:"outDir"`
	Base         string              `json:"base"`
	Output       string              `json:"output"`
	Build        BuildConfig         `json:"build"`
	Integrations []IntegrationConfig `json:"integrations"`
}

func OpenConfig(projectDir string) (*ConfigInfo, error) {
	cmd := exec.Command("node", "-e", openConfigSrc)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var c ConfigInfo
	json.Unmarshal(output, &c)

	return &c, nil
}

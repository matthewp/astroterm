package build

import (
	"astroterm/project"
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed astroterm-build.config.mjs.tmpl
var astroTermBuildConfig string

type templateData struct {
	Path string
}

type ConfigBuilder struct {
	tmpl *template.Template
}

func NewConfigBuilder() *ConfigBuilder {
	template, _ := template.New("config").Parse(astroTermBuildConfig)
	return &ConfigBuilder{
		tmpl: template,
	}
}

func (b *ConfigBuilder) CreateBuildConfig(p *project.Project) ([]byte, error) {
	data := &templateData{
		Path: p.ConfigPath(),
	}

	var tpl bytes.Buffer
	err := b.tmpl.Execute(&tpl, data)
	if err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}

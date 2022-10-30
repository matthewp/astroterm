package project

import (
	aenv "astroterm/env"
)

type Project struct {
	Dir string
	env *aenv.TermEnvironment
	pkg *PackageJson
}

func NewProject() *Project {
	return &Project{
		env: nil,
		pkg: nil,
	}
}

func OpenLocalProject() (*Project, error) {
	env, err := aenv.GetEnvironment()
	if err != nil {
		return nil, err
	}

	proj := NewProject()
	proj.env = env
	proj.Dir = env.Pwd

	pkg, _ := OpenPackageJson(env)
	proj.pkg = pkg
	return proj, nil
}

func (p *Project) Name() string {
	if p.pkg == nil {
		return "unknown"
	}
	return p.pkg.Name
}

func (p *Project) PackageJson() *PackageJson {
	return p.pkg
}

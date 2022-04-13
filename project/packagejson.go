package project

import (
	aenv "astroterm/env"
	"encoding/json"
	"io/ioutil"
	"path"
)

type PackageJson struct {
	name string
	deps map[string]string
}

func OpenPackageJson(env *aenv.TermEnvironment) (*PackageJson, error) {
	pkgPath := path.Join(env.Pwd, "package.json")
	bytes, err := ioutil.ReadFile(pkgPath)

	if err != nil {
		return nil, err
	}

	pkg := &PackageJson{}
	err = json.Unmarshal(bytes, pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

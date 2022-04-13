package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DevServerModel struct {
	Pid        int
	Hostname   string
	Port       int
	ProjectDir string
}

func (d *Database) LoadDeverServerModel(projectDir string) (*DevServerModel, error) {
	if err := d.ensureOpened(); err != nil {
		return nil, err
	}
	r, err := d.db.Exec("SELECT * FROM devservers WHERE projectdir = ?;", projectDir)

	fmt.Printf("%v", r)

	return nil, err
}

func (d *Database) AddStartedDevServer(model *DevServerModel) error {
	if err := d.ensureOpened(); err != nil {
		return err
	}
	_, err := d.db.Exec("INSERT INTO devservers VALUES(?,NULL,NULL,?);", model.Pid, model.ProjectDir)
	return err
}

func (d *Database) SetDevServerInformation(model *DevServerModel) error {
	if err := d.ensureOpened(); err != nil {
		return err
	}
	_, err := d.db.Exec(`UPDATE devservers
	SET port = ?, hostname = ?
	WHERE pid = ?;`, model.Port, model.Hostname, model.Pid)
	return err
}

func (d *Database) DeleteDevServer(model *DevServerModel) error {
	if err := d.ensureOpened(); err != nil {
		return err
	}
	_, err := d.db.Exec(`DELETE FROM devservers WHERE pid = ?;`, model.Pid)
	return err
}

package db

import (
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
	stmt, err := d.db.Prepare("SELECT pid, hostname, port, projectdir FROM devservers WHERE projectdir = ?;")
	if err != nil {
	}
	defer stmt.Close()

	model := &DevServerModel{}
	row := stmt.QueryRow(projectDir)
	err = row.Scan(&model.Pid, &model.Hostname, &model.Port, &model.ProjectDir)
	if model.Pid == 0 {
		return nil, err
	}
	return model, err
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

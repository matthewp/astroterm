package db

import (
	_ "github.com/mattn/go-sqlite3"
)

type DevServerModel struct {
	Pid        int
	Hostname   string
	Port       int
	Subpath    string
	ProjectDir string
	LogPath    string
}

func (m *DevServerModel) IsRunning() bool {
	return m.Pid != 0
}

func (d *Database) LoadDevServerModel(projectDir string) (*DevServerModel, error) {
	if err := d.ensureOpened(); err != nil {
		return nil, err
	}
	stmt, err := d.db.Prepare("SELECT pid, hostname, port, subpath, projectdir, logpth FROM devservers WHERE projectdir = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	model := &DevServerModel{}
	row := stmt.QueryRow(projectDir)
	err = row.Scan(&model.Pid, &model.Hostname, &model.Port, &model.Subpath, &model.ProjectDir, &model.LogPath)
	if model.Pid == 0 {
		return nil, err
	}
	return model, err
}

func (d *Database) AddStartedDevServer(model *DevServerModel) error {
	if err := d.ensureOpened(); err != nil {
		return err
	}
	_, err := d.db.Exec(`INSERT INTO devservers VALUES(?,0,"","",?,?);`, model.Pid, model.ProjectDir, model.LogPath)
	return err
}

func (d *Database) SetDevServerInformation(model *DevServerModel) error {
	if err := d.ensureOpened(); err != nil {
		return err
	}
	_, err := d.db.Exec(`UPDATE devservers
	SET port = ?, hostname = ?, subpath = ?
	WHERE pid = ?;`, model.Port, model.Hostname, model.Subpath, model.Pid)
	return err
}

func (d *Database) DeleteDevServer(model *DevServerModel) error {
	if err := d.ensureOpened(); err != nil {
		return err
	}
	_, err := d.db.Exec(`DELETE FROM devservers WHERE pid = ?;`, model.Pid)
	return err
}

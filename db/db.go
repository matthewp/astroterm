package db

import (
	"astroterm/home"
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "stateV0.db"
const createSchemaStatement string = `
CREATE TABLE IF NOT EXISTS devservers (
	pid 		integer NOT NULL PRIMARY KEY,
	port 		integer,
	hostname 	text,
	subpath 	text,
	projectdir 	text,
	logpth		text
)`

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	return &Database{
		db: nil,
	}
}

func (d *Database) Open() error {
	homeFolder, err := home.OpenHomeFolder()
	if err != nil {
		return err
	}

	dbPath := path.Join(homeFolder, dbName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	d.db = db
	if err = d.createSchema(); err != nil {
		return err
	}
	return nil
}

func (d *Database) createSchema() error {
	_, err := d.db.Exec(createSchemaStatement)
	return err
}

func (d *Database) ensureOpened() error {
	if d.db == nil {
		return d.Open()
	}
	return nil
}

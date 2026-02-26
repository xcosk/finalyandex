package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(256) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date_idx ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	if dbFile == "" {
		return errors.New("db file is empty")
	}

	_, err := os.Stat(dbFile)
	install := err != nil
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			db.Close()
			return err
		}
	}

	DB = db
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	return DB.Close()
}

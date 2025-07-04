package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB

const schema = `CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(100) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(100) NOT NULL DEFAULT ""
		);`

const indexSchema = "CREATE INDEX idx_date ON scheduler(date);"

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	dbCon, err := sql.Open("sqlite", dbFile)

	if install {
		if _, err := dbCon.Exec(schema); err != nil {
			return err
		}
		if _, err := dbCon.Exec(indexSchema); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	db = dbCon

	return nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

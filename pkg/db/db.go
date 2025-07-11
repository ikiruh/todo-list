package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	if install {
		err := createTableScheduler()
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func createTableScheduler() error {
	scheduler := `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(30) NOT NULL DEFAULT "",
		comment TEXT,
		repeat VARCHAR(128)
	);
	CREATE INDEX IF NOT EXISTS idx_schedule_title ON scheduler(title);`

	_, err := db.Exec(scheduler)
	return err
}

package db

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

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

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := db.NamedExec(query, task)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetTasks(limit int, search string) ([]*Task, error) {
	var tasks []*Task
	var err error

	if date, err := time.Parse("02.01.2006", search); err == nil {
		formattedDate := date.Format("20060102")
		query := `SELECT id, date, title, comment, repeat FROM scheduler 
					WHERE date = ? ORDER BY date ASC LIMIT ?`
		err = db.Select(&tasks, query, formattedDate, limit)
		if err != nil {
			return []*Task{}, err
		}
	} else if search != "" {
		searchPattern := "%" + search + "%"
		query := `SELECT id, date, title, comment, repeat FROM scheduler 
					WHERE title LIKE ? OR comment LIKE ? 
					ORDER BY date ASC LIMIT ?`
		err = db.Select(&tasks, query, searchPattern, searchPattern, limit)
		if err != nil {
			return []*Task{}, err
		}
	} else {
		query := `SELECT id, date, title, comment, repeat FROM scheduler 
					ORDER BY date ASC LIMIT ?`
		err = db.Select(&tasks, query, limit)
		if err != nil {
			return []*Task{}, err
		}
	}

	if err != nil {
		return []*Task{}, err
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

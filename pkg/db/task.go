package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
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

func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.Get(&task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task mot found")
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
				 SET date = :date, title = :title, comment = :comment, repeat = :repeat
				 WHERE id = :id`

	result, err := db.NamedExec(query, task)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
}

func DeleteTask(id int64) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateDate(id int64, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, newDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

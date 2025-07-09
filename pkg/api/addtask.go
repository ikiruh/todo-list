package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ikiruh/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Error decode JSON", http.StatusBadRequest)
		return
	}
	if err := validateTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "Error add task: %w", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"id": id}, http.StatusCreated)
}

func validateTask(task *db.Task) error {
	if task.Title == "" {
		return fmt.Errorf("the issue title is not specified")
	}

	now := time.Now()
	today := now.Format(LayoutDate)

	if task.Date == "" {
		task.Date = today
	} else {
		t, err := time.Parse(LayoutDate, task.Date)
		if err != nil {
			return fmt.Errorf("incorrect date format")
		}

		if afterNow(now, t) {
			if task.Repeat == "" {
				task.Date = today
			} else {
				next, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return err
				}
				task.Date = next
			}
		}
	}

	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			return err
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	writeJSON(w, map[string]string{"error": message}, statusCode)
}

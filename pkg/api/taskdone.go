package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ikiruh/go_final_project/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "task ID is not specified", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(idStr)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	now := time.Now()

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.UpdateDate(id, nextDate); err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}

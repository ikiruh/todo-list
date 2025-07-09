package api

import (
	"net/http"
	"strconv"

	"github.com/ikiruh/go_final_project/pkg/db"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "The method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			writeJSONError(w, "Invalid limit value", http.StatusBadRequest)
			return
		}
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.GetTasks(limit, search)
	if err != nil {
		writeJSONError(w, "Error when receiving tasks", http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResponse{
		Tasks: tasks,
	}, http.StatusOK)
}

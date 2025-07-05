package api

import (
	"encoding/json"
	"final-project/pkg/db"
	"io"
	"net/http"
	"strconv"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var id int

	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "некорректный ID", http.StatusBadRequest)
			return
		}
	}

	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id == 0 {
		id = task.ID
	}

	if id <= 0 {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	task.ID = id

	err = db.UpdateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

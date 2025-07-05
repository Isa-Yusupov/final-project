package api

import (
	"final-project/pkg/db"
	"net/http"
	"strconv"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, task)
}

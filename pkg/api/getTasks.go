package api

import (
	"final-project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, _ *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		http.Error(w, "ошибка json", http.StatusBadRequest)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

package api

import (
	"final-project/pkg/db"
	"fmt"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func messageError(err error) string {
	return fmt.Sprintf(`{"error":"%s"}`, err.Error())
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error
	if searchQuery != "" {
		tasks, err = db.GetTasks(50, searchQuery)
	} else {
		tasks, err = db.GetTasks(50, "") // в параметре максимальное количество записей
	}

	if err != nil {
		http.Error(w, messageError(err), http.StatusBadRequest)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

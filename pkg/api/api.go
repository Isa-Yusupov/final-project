package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
}

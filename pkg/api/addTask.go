package api

import (
	"encoding/json"
	"final-project/pkg/db"
	"io"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "неудалось десериализировать json", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "title обязательно", http.StatusBadRequest)
		return
	}
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ошибка с датой"))
			return
		}
		task.Date = next
	}

	taskTime, err := time.Parse(dateFormat, task.Date)
	if err != nil || taskTime.Before(now) {
		task.Date = now.Format(dateFormat)
	}

	taskID, err := db.AddTask(&task)
	if err != nil {
		http.Error(w, "не удалось сохранить задачу", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, taskID)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "не удалось отправить ответ", http.StatusInternalServerError)
	}
}

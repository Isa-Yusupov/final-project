package api

import (
	"final-project/pkg/db"
	"net/http"
	"time"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	task, err := db.GetTask(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		return
	}

	now, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	nextD, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	err = db.UpdateDate(nextD, idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

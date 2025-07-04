package api

import (
	"final-project/pkg/db"
	"net/http"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	db.DeleteTask(idStr)
}

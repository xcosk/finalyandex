package api

import (
	"net/http"

	"finalyandex/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasksResp{Tasks: tasks})
}

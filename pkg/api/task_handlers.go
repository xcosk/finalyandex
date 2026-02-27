package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"finalyandex/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.FormValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	t, err := db.GetTask(strconv.FormatInt(id, 10))
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req db.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	_, err := parseID(req.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is empty")
		return
	}
	if err := validateRepeat(req.Repeat); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	date, err := normalizeTaskDate(time.Now(), req.Date, req.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Date = date

	if err := db.UpdateTask(&req); err != nil {
		if strings.Contains(err.Error(), "incorrect id") {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.FormValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.DeleteTask(strconv.FormatInt(id, 10)); err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.FormValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sid := strconv.FormatInt(id, 10)
	task, err := db.GetTask(sid)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(sid); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateDate(next, sid); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

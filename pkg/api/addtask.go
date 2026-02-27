package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"finalyandex/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(dateLayout)
	}

	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return err
	}

	next := ""
	if strings.TrimSpace(task.Repeat) != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "title is empty")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func afterNow(date, now time.Time) bool {
	return dateOnly(date).After(dateOnly(now))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat")
		}
		date = date.AddDate(1, 0, 0)
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(dateLayout), nil
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errors.New("invalid repeat")
		}
		date = date.AddDate(0, 0, interval)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(dateLayout), nil
	default:
		return "", errors.New("unsupported repeat")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	now := time.Now()
	if nowStr != "" {
		parsed, err := time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(next))
}

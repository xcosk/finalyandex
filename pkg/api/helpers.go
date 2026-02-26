package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

func parseID(raw string) (int64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, errors.New("id is empty")
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func normalizeTaskDate(now time.Time, date, repeat string) (string, error) {
	today := now.Format(dateLayout)
	if strings.TrimSpace(date) == "" {
		return today, nil
	}

	parsed, err := time.Parse(dateLayout, date)
	if err != nil {
		return "", err
	}

	if afterNow(parsed, now) || dateOnly(parsed).Equal(dateOnly(now)) {
		return date, nil
	}

	if strings.TrimSpace(repeat) == "" {
		return today, nil
	}
	return NextDate(now, date, repeat)
}

func validateRepeat(repeat string) error {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return nil
	}
	_, err := NextDate(time.Now(), time.Now().Format(dateLayout), repeat)
	return err
}

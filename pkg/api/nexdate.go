package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	splitRepeat := strings.Split(repeat, " ")
	var interval int
	var value string

	if len(splitRepeat) == 1 {
		value = splitRepeat[0]
		interval = 1
	} else {
		var err error
		value = splitRepeat[0]
		interval, err = strconv.Atoi(splitRepeat[1])
		if err != nil {
			return "", err
		}
	}
	date, err := time.Parse(dateFormat, dstart)

	if err != nil {
		return "", err
	}

	if afterNow(date, now) {
		return "", fmt.Errorf("начальная дата должна быть в прошлом или равна now")
	}

	if err != nil {
		return "", fmt.Errorf("неверный интервал: %w", err)
	}

	switch value {
	case "d":
		if interval > 400 {
			interval = 400
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}
	case "y":
		for {
			date = date.AddDate(interval, 0, 0)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}
	default:
		return "", fmt.Errorf("repeat не поддерживается: %s")
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func nextDayHandler(w http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeat := req.FormValue("repeat")

	if dateStr == "" || repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	var now time.Time
	var err error

	if nowStr != "" {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "неверный формат даты now", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}

	resp, err := NextDate(now, dateStr, repeat)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server Error"))
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

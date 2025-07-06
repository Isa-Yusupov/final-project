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

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	value := splitRepeat[0]

	switch value {
	case "d":
		if len(splitRepeat) != 2 {
			return "", errors.New("формат d требует интервала")
		}

		interval, err := strconv.Atoi(splitRepeat[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("неверный интервал для d (1-400)")
		}

		if afterNow(date, now) {
			return date.Format(dateFormat), nil
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "y":
		if len(splitRepeat) != 1 {
			return "", errors.New("формат y не требует параметров")
		}

		if afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "w":
		if len(splitRepeat) != 2 {
			return "", errors.New("формат w требует дней недели")
		}
		dayMap := map[int]bool{}
		days := strings.Split(splitRepeat[1], ",")
		for _, d := range days {
			day, err := strconv.Atoi(d)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("неверный день недели: %s", d)
			}
			dayMap[day] = true
		}
		// искать следующую дату после now, которая совпадает с одним из дней недели
		current := now.AddDate(0, 0, 1)
		for {
			weekday := int(current.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if dayMap[weekday] {
				return current.Format(dateFormat), nil
			}
			current = current.AddDate(0, 0, 1)
		}
	case "m":
		if len(splitRepeat) < 2 {
			return "", errors.New("формат m требует минимум дни месяца")
		}
		days := strings.Split(splitRepeat[1], ",")
		months := map[int]bool{}
		dayList := []int{}

		for _, d := range days {
			day, err := strconv.Atoi(d)
			if err != nil || (day == 0 || (day < -2 || day > 31)) {
				return "", fmt.Errorf("неверный день месяца: %s", d)
			}
			dayList = append(dayList, day)
		}

		if len(splitRepeat) > 2 {
			monthParts := strings.Split(splitRepeat[2], ",")
			for _, m := range monthParts {
				mon, err := strconv.Atoi(m)
				if err != nil || mon < 1 || mon > 12 {
					return "", fmt.Errorf("неверный месяц: %s", m)
				}
				months[mon] = true
			}
		}

		// искать следующий подходящий день
		current := now.AddDate(0, 0, 1)
		for {
			if len(months) > 0 && !months[int(current.Month())] {
				current = current.AddDate(0, 0, 1)
				continue
			}

			lastDay := time.Date(current.Year(), current.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			for _, d := range dayList {
				targetDay := d
				if d < 0 {
					targetDay = lastDay + d + 1
				}
				if targetDay < 1 || targetDay > lastDay {
					continue
				}
				if current.Day() == targetDay {
					return current.Format(dateFormat), nil
				}
			}

			current = current.AddDate(0, 0, 1)
		}
	default:
		return "", fmt.Errorf("repeat не поддерживается: %s")
	}
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC).After(time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC))
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

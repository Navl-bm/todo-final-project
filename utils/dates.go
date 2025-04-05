package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/Navl-bm/todo-final-project/errors"
)

const DateFormat = "20060102"

func AfterNow(a, b time.Time) bool {
	return a.Compare(b) == 1
}

func getMonthOfYear(month int) (time.Month, bool) {
	switch month {
	case 1:
		return time.January, true
	case 2:
		return time.February, true
	case 3:
		return time.March, true
	case 4:
		return time.April, true
	case 5:
		return time.May, true
	case 6:
		return time.June, true
	case 7:
		return time.July, true
	case 8:
		return time.August, true
	case 9:
		return time.September, true
	case 10:
		return time.October, true
	case 11:
		return time.November, true
	case 12:
		return time.December, true
	}
	return 0, false
}

func getDayOfWeek(day int) (time.Weekday, bool) {
	switch day {
	case 1:
		return time.Monday, true
	case 2:
		return time.Tuesday, true
	case 3:
		return time.Wednesday, true
	case 4:
		return time.Thursday, true
	case 5:
		return time.Friday, true
	case 6:
		return time.Saturday, true
	case 7:
		return time.Sunday, true
	}
	return 0, false
}

func nextDay(parts []string, date time.Time, now time.Time) (string, error) {
	if len(parts) != 2 {
		return "", errors.ErrBadRepeatFormat
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.ErrBadRepeatFormat
	}
	if interval > 400 {
		return "", errors.ErrBadRepeatFormat
	}
	for {
		date = date.AddDate(0, 0, interval)
		if AfterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func nextWeek(parts []string, date time.Time, now time.Time) (string, error) {
	if len(parts) != 2 {
		return "", errors.ErrBadRepeatFormat
	}

	var days []time.Weekday
	for _, v := range strings.Split(parts[1], ",") {
		val, err := strconv.Atoi(v)
		if err != nil {
			return "", errors.ErrBadRepeatFormat
		}

		day, ok := getDayOfWeek(val)
		if !ok {
			return "", errors.ErrBadRepeatFormat
		}

		days = append(days, day)
	}
	for {
		date = date.AddDate(0, 0, 1)

		if AfterNow(date, now) {
			for _, v := range days {
				if date.Weekday() == v {
					return date.Format(DateFormat), nil
				}
			}
		}
	}
}

func nextMonth(parts []string, date time.Time, now time.Time) (string, error) {
	if len(parts) > 3 || len(parts) == 1 {
		return "", errors.ErrBadRepeatFormat
	}

	daysString := strings.Split(parts[1], ",")
	var days []int
	for _, v := range daysString {
		tmp, err := strconv.Atoi(v)
		if err != nil {
			return "", errors.ErrBadRepeatFormat
		}
		if tmp < -2 || tmp > 31 {
			return "", errors.ErrBadRepeatFormat
		}
		days = append(days, tmp)
	}

	var monthString []string
	var month []time.Month
	monthExists := false
	if len(parts) == 3 {
		monthExists = true
		monthString = strings.Split(parts[2], ",")
		for _, v := range monthString {
			tmp, err := strconv.Atoi(v)
			if err != nil {
				return "", errors.ErrBadRepeatFormat
			}
			if tmp > 12 || tmp < 1 {
				return "", errors.ErrBadRepeatFormat
			}

			monthName, ok := getMonthOfYear(tmp)
			if !ok {
				return "", errors.ErrBadRepeatFormat
			}
			month = append(month, monthName)
		}
	}

	date = date.AddDate(0, 0, 1)

	if monthExists {
		for {
			if AfterNow(date, now) {
				for _, m := range month {
					if date.Month() == m {
						for {
							for _, v := range days {
								if v == -1 && date.AddDate(0, 0, 1).Month() == time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC).Month() {
									return date.Format(DateFormat), nil
								} else if v == -2 && date.AddDate(0, 0, 2).Month() == time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC).Month() {
									return date.Format(DateFormat), nil
								} else if date.Day() == v {
									return date.Format(DateFormat), nil
								}
							}
							date = date.AddDate(0, 0, 1)
						}
					}
				}
			}
			date = time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		}
	} else {
		for {
			if AfterNow(date, now) {
				for _, v := range days {
					if v == -2 && date.AddDate(0, 0, 2).Month() == time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC).Month() {
						return date.Format(DateFormat), nil
					} else if v == -1 && date.AddDate(0, 0, 1).Month() == time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC).Month() {
						return date.Format(DateFormat), nil
					} else if date.Day() == v {
						return date.Format(DateFormat), nil
					}
				}
			}
			date = date.AddDate(0, 0, 1)
		}
	}
}

func nextYear(parts []string, date time.Time, now time.Time) (string, error) {
	if len(parts) != 1 {
		return "", errors.ErrBadRepeatFormat
	}
	for {
		date = date.AddDate(1, 0, 0)
		if AfterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func NextDate(dnow string, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.ErrBadRepeatFormat
	}
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.ErrBadTime
	}
	now, err := time.Parse(DateFormat, dnow)
	if err != nil {
		return "", errors.ErrBadTime
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		{
			return nextDay(parts, date, now)
		}
	case "w":
		{
			return nextWeek(parts, date, now)
		}
	case "m":
		{
			return nextMonth(parts, date, now)
		}
	case "y":
		{
			return nextYear(parts, date, now)
		}
	default:
		{
			return "", errors.ErrBadRepeatFormat
		}
	}
}

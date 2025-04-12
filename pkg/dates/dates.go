package dates

import (
	"strconv"
	"strings"
	"time"

	"github.com/Navl-bm/todo-final-project/errors"
)

const DateFormat = "20060102"

// getMonthOfYear возвращает значение типа time.Month
// в зависимости от переданного числового значения месяца
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

// getDayOfWeek возвращает значение типа time.Weekday
// в зависимости от переданного числового значения дня недели
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

// nextDay возвращает ближайшую дату при выбранном повторении по дням
// формат d <число>
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

		if date.After(now) {
			return date.Format(DateFormat), nil
		}
	}
}

// nextWeek возвращает дату при выбранном повторении по неделям
// формат w <через запятую от 1 до 7>
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

		if date.After(now) {
			for _, v := range days {
				if date.Weekday() == v {
					return date.Format(DateFormat), nil
				}
			}
		}
	}
}

// findNextMonthDay вспомогательная функция для поиска следующего месяца и дня
func findNextMonthDay(monthExists bool, date, now time.Time, months []time.Month, days []int) (string, error) {
	// получаем дату и возвращаем первый день следующего месяца
	nextMonthStart := func(date time.Time) time.Time {
		return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	}

	// проверка даты на соответствие условию повторения по дню
	isDayMatch := func(date time.Time, day int) bool {
		switch day {
		case -1:
			return date.AddDate(0, 0, 1).Month() == nextMonthStart(date).Month()
		case -2:
			return date.AddDate(0, 0, 2).Month() == nextMonthStart(date).Month()
		default:
			return date.Day() == day
		}
	}

	// проверка даты на соответствие условию повторения по месяцу
	checkMonth := func(date time.Time) bool {
		if !monthExists {
			return true
		}
		for _, m := range months {
			if date.Month() == m {
				return true
			}
		}
		return false
	}

	for {
		// пока дата не превысит(станет равной текущему месяцу) к текущей дате прибавляем месяц
		if !date.After(now) {
			date = nextMonthStart(date)
			continue
		}

		// проверка месяца на соответствие условию
		if checkMonth(date) {
			currentDate := date
			for {
				// проходим по всем дням условия и проверяем на соответствие
				for _, day := range days {
					if isDayMatch(currentDate, day) {
						return currentDate.Format(DateFormat), nil
					}
				}
				currentDate = currentDate.AddDate(0, 0, 1)

				// если дата перешла в следующий месяц выходим из цикла
				if currentDate.Month() != date.Month() {
					break
				}
			}
		}

		// переходим к следующему месяцу, т.к. условия не выполнены
		date = nextMonthStart(date)
	}
}

// nextMonth возвращает ближайщую дату при повторении по месяцам
// формат m <через запятую от 1 до 31, -1, -2> [через запятую от 1 до 12]
// (-1) - последний день месяца, (-2) - предпоследний день месяца
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

	nextDate, err := findNextMonthDay(monthExists, date, now, month, days)
	if err != nil {
		return "", err
	}

	return nextDate, nil
}

// nextYear возвращает ближайщую дату при повторении по годам
// формат y
func nextYear(parts []string, date time.Time, now time.Time) (string, error) {
	if len(parts) != 1 {
		return "", errors.ErrBadRepeatFormat
	}
	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			return date.Format(DateFormat), nil
		}
	}
}

// NextDate определяет какому типу соответствует правило повторения
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

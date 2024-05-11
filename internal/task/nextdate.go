package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	validDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("incorrect date %v", err)
	}

	rule := string(repeat[0])
	rightLen := len(repeat) > 2
	var result string

	switch {
	//задача переносится на указанное число дней
	case rule == "d" && rightLen:
		result, err = everyDay(now, validDate, repeat[2:])

	// задача назначается в указанные дни недели, где 1 — понедельник, 7 — воскресенье
	case rule == "w" && rightLen:
		result, err = everyWeek(validDate, now, repeat[2:])

	// задача назначается в указанные дни месяца (1-31)
	case rule == "m" && rightLen:
		result, err = everyMonth(validDate, now, repeat[2:])

	// задача выполняется ежегодно
	case rule == "y":
		result, err = everyYear(now, validDate)
	default:
		return "", fmt.Errorf("incorrect repetition rule %v", err)
	}

	return result, err
}

func everyDay(now, date time.Time, days string) (string, error) {
	d, err := strconv.Atoi(days)
	if err != nil || d > 400 || d < 0 {
		return "", fmt.Errorf(`incorrect repetition rule in "d"`)
	}

	resultDate := date.AddDate(0, 0, d)
	for resultDate.Before(now) {
		resultDate = resultDate.AddDate(0, 0, d)
	}

	return resultDate.Format("20060102"), nil
}

func everyWeek(date, now time.Time, repeat string) (string, error) {
	result := ""

	week := make(map[int]string)

	if date.Before(now) {
		date = now
	}

	days := strings.Split(string(repeat), ",")

	for i := 1; i <= 7; i++ {
		date = date.AddDate(0, 0, 1)
		weekDay := int(date.Weekday())

		if weekDay == 0 {
			weekDay = 7
		}

		week[weekDay] = date.Format("20060102")

		for _, day := range days {
			d, err := strconv.Atoi(day)
			if err != nil || d > 7 || d < 0 {
				return "", fmt.Errorf(`incorrect repetition rule in "w" %v`, err)
			}

			if d == weekDay {
				result = week[d]
				return result, nil
			}
		}
	}

	return result, nil
}

func everyMonth(date, now time.Time, repeat string) (string, error) {
	return "", nil
}

// функция для получения последнего дня месяца
func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

func everyYear(now, date time.Time) (string, error) {
	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(1, 0, 0)
		}
	} else {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format("20060102"), nil
}

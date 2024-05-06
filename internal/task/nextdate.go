package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" || repeat == "d 1" {
		return now.Format("20060102"), nil
	}

	realDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("incorrect date %v", err)
	}

	rule := string(repeat[0])
	rightLen := len(repeat) > 2
	var result string

	switch {
	//задача переносится на указанное число дней
	case rule == "d" && rightLen:
		result, err = everyDay(now, realDate, repeat[2:])

	// задача назначается в указанные дни недели, где 1 — понедельник, 7 — воскресенье
	case rule == "w" && rightLen:
		result, err = everyWeek(realDate, now, repeat[2:])

	// задача назначается в указанные дни месяца (1-31)
	case rule == "m" && rightLen:
		result, err = everyMonth(realDate, now, repeat[2:])

	// задача выполняется ежегодно
	case rule == "y":
		result, err = everyYear(now, realDate)
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

	if d == 1 && date.After(now) {
		return date.AddDate(0, 0, 1).Format("20060102"), nil
	}

	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(0, 0, d)
		}
	} else {
		date = date.AddDate(0, 0, d)
	}

	return date.Format("20060102"), nil
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
	result := ""

	nextYear := make(map[int]string)

	// если дата раньше текущей, то меняем ее на текущую
	if date.Before(now) {
		date = now
	}

	// получаем количество аргументов repeat
	args := strings.Split(repeat, " ")

	// первый аргумент - по каким дням
	days := strings.Split(args[0], ",")

	// второй аргумент - по каким месяцам
	months := make([]string, 0)
	if len(args) > 1 {
		months = strings.Split(args[1], ",")
	}

	for i := 1; i <= 365; i++ {
		// переходим к следующему дню
		date = date.AddDate(0, 0, 1)

		// определяем месяц
		monthDay := int(date.Day())

		// добавляем следующий день в мапу
		nextYear[monthDay] = date.Format("20060102")

		// итерируемся по дням
		for _, day := range days {
			d, err := strconv.Atoi(day)
			if err != nil || d > 31 {
				return "", fmt.Errorf(`incorrect repetition rule in "m" %v`, err)
			}

			// если не указаны месяцы и день - отрицательное число
			if d < 0 && len(months) == 0 {
				result := EndOfMonth(date)
				if d == -1 {
					return result.Format("20060102"), nil
				}
				result = result.AddDate(0, 0, d-1)
				return result.Format("20060102"), nil
			}

			// если не указаны месяцы и текущий день совпадает с нужным днем
			if d == monthDay && len(months) == 0 {
				result = nextYear[d]
				return result, nil
			}

			// если указаны месяцы и текущий день совпадает с нужным днем
			if d == monthDay && len(months) != 0 {

				for _, month := range months {
					m, err := strconv.Atoi(month)
					if err != nil || m < 0 || m > 12 {
						return "", fmt.Errorf(`incorrect repetition rule in "m" %v`, err)
					}

					dateMonth, _ := time.Parse("20060102", nextYear[d])
					currMonth := int(dateMonth.Month())

					if m == currMonth {
						result = nextYear[d]
						return result, nil
					} else {
						continue
					}
				}
			}

		}

	}

	return result, nil
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

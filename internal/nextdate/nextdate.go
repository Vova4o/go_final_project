package nextdate

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate returns the next date of the task
// с параметрами:
// now — время от которого ищется ближайшая дата;
// date — исходное время в формате 20060102, от которого начинается отсчёт повторений;
// repeat — правило повторения в описанном выше формате.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch repeat[0] {
	case 'y':
		return addYear(t, now)
	case 'd':
		return addDays(t, now, repeat)
	case 'w':
		// err := errors.New("не реализовано")
		// return date, err
		return AddWeeks(t, now, repeat)
	case 'm':
		// err := errors.New("не реализовано")
		// return date, err
		return AddMonths(t, now, repeat)
	}
	return "", nil
}

func addYear(t time.Time, now time.Time) (string, error) {
	for {
		t = t.AddDate(1, 0, 0)
		if t.After(now) {
			break
		}
	}
	return t.Format("20060102"), nil
}

func addDays(t time.Time, now time.Time, repeat string) (string, error) {
	rep := strings.Split(repeat, " ")
	if len(rep) != 2 {
		err := errors.New("не указан интервал в днях")
		return "", err
	}
	daysNumber, err := strconv.Atoi(rep[1])
	if err != nil {
		return "", err
	}
	if daysNumber > 400 || daysNumber < 1 {
		err := errors.New("превышен максимально допустимый интервал")
		return "", err
	} else {
		for {
			t = t.AddDate(0, 0, daysNumber)
			if t.After(now) {
				break
			}
		}
		return t.Format("20060102"), nil
	}
}

func AddWeeks(t time.Time, now time.Time, repeat string) (string, error) {
	var weekDays []int
	if len(repeat) < 3 {
		return "", errors.New("не указан интервал в днях недели")
	}

	weekDaysStr := strings.Split(repeat[2:], ",")
	if len(weekDaysStr) == 0 {
		return "", errors.New("не указан интервал в днях недели")
	}

	for _, day := range weekDaysStr {
		dayNumber, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}
		if dayNumber < 1 || dayNumber > 7 {
			return "", fmt.Errorf("недопустимое значение %d дня недели", dayNumber)
		}
		weekDays = append(weekDays, dayNumber)
	}

	for i, day := range weekDays {
		if day == 7 {
			weekDays[i] = 0
		}
	}

	sort.Ints(weekDays)
	var nextWeekDay int
	for _, wd := range weekDays {
		if wd >= int(t.Weekday()) { // Check if wd is greater than or equal to the current weekday
			nextWeekDay = wd
			break
		}
	}
	if nextWeekDay == 0 { // If no future weekday was found in this week, take the first day of the next week
		nextWeekDay = weekDays[0]
	}
	for {
		t = t.AddDate(0, 0, 1)
		if t.After(now) && int(t.Weekday()) == nextWeekDay {
			return t.Format("20060102"), nil
		}
	}
}

func AddMonths(t time.Time, now time.Time, repeat string) (string, error) {
	var listOfDays time.Time
	var err error
	if len(repeat) < 3 {
		return "", errors.New("не указан интервал в месяцах")
	}

	repSlice := strings.Split(repeat, " ")

	if len(repSlice) == 2 {
		listOfDays, err = getNextDate(now, t.Format("20060102"), repeat)
		if err != nil {
			return "", err
		}
	} else if len(repSlice) == 3 {
		listOfDays = getNextMonthDate(now, t.Format("20060102"), repeat)
	} else {
		return "", errors.New("неверный формат повторения")
	}

	return listOfDays.Format("20060102"), nil
}

func getNextMonthDate(now time.Time, target string, rule string) time.Time {
	// Parse the rule
	ruleParts := strings.Split(rule, " ")
	daysPart := strings.Split(ruleParts[1], ",")
	monthsPart := strings.Split(ruleParts[2], ",")

	// Convert days and months to integers
	days := make([]int, len(daysPart))
	for i, day := range daysPart {
		days[i], _ = strconv.Atoi(day)
	}
	months := make([]int, len(monthsPart))
	for i, month := range monthsPart {
		months[i], _ = strconv.Atoi(month)
	}

	// Initialize targetTime and nearestDate
	targetTime, _ := time.Parse("20060102", target)
	var nearestDate time.Time

	// Loop until we find a date after today
	for {
		for _, day := range days {
			for _, month := range months {
				var date time.Time
				if day < 0 {
					// If day is negative, calculate the date from the end of the current month
					endOfMonth := time.Date(targetTime.Year(), time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
					date = endOfMonth.AddDate(0, 0, day+1)
				} else {
					date = time.Date(targetTime.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
				}
				// If the date is before now, add a year
				if date.Before(now) {
					date = date.AddDate(1, 0, 0)
				}
				// If this is the first date we've found, or it's earlier than the current nearest date, update nearestDate
				if (nearestDate.IsZero() || date.Before(nearestDate)) && date.After(now) {
					nearestDate = date
				}
			}
		}
		// If nearestDate is not zero, it means we have found a date after today, so break the loop
		if !nearestDate.IsZero() {
			break
		}
		// If no date after today is found in this year, increment the year
		targetTime = targetTime.AddDate(1, 0, 0)
	}

	return nearestDate
}

func getNextDate(now time.Time, target string, rule string) (time.Time, error) {
	// Parse target date
	targetTime, _ := time.Parse("20060102", target)

	// Split rule into "m" and the rest
	ruleParts := strings.SplitN(rule, " ", 2)
	if len(ruleParts) != 2 {
		return time.Time{}, fmt.Errorf("invalid rule format")
	}

	// Check if the rest of the rule contains a comma
	var monthDays []int
	if strings.Contains(ruleParts[1], ",") {
		daysParts := strings.Split(ruleParts[1], ",")
		monthDays = make([]int, len(daysParts))
		for i, part := range daysParts {
			day, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid day in rule: %v", err)
			}
			if day < -2 || day > 31 {
				return time.Time{}, fmt.Errorf("не правильно указан формат повтора")
			}
			monthDays[i] = day
		}
	} else {
		day, err := strconv.Atoi(strings.TrimSpace(ruleParts[1]))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day in rule: %v", err)
		}
		if day < -2 || day > 31 {
			return time.Time{}, fmt.Errorf("не правильно указан формат повтора")
		}
		monthDays = []int{day}
	}
	// Find the nearest date
	var nearestDate time.Time
	for {
		for _, day := range monthDays {
			var date time.Time
			// Define endOfMonth here
			endOfMonth := time.Date(targetTime.Year(), targetTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			if day < 0 {
				// If day is negative, calculate the date from the end of the current month
				date = endOfMonth.AddDate(0, 0, day+1)
			} else {
				if day > endOfMonth.Day() {
					targetTime = targetTime.AddDate(0, 1, 0)
					endOfMonth = time.Date(targetTime.Year(), targetTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)
				}
				date = time.Date(targetTime.Year(), targetTime.Month(), day, 0, 0, 0, 0, time.UTC)
			}
			// Now endOfMonth is accessible here
			if date.Before(now) || date.After(endOfMonth) {
				date = date.AddDate(0, 1, 0)
			}
			// If this is the first date we've found, or it's earlier than the current nearest date, update nearestDate
			if (nearestDate.IsZero() || date.Before(nearestDate)) && date.After(now) {
				nearestDate = date
				// break
			}
		}
		// If nearestDate is not zero, it means we have found a date after today, so break the loop
		if !nearestDate.IsZero() {
			break
		}
		// If no date after today is found in this month, increment the month
		targetTime = targetTime.AddDate(0, 1, 0)
	}

	return nearestDate, nil
}

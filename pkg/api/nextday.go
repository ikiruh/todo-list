package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}
	date, err := time.Parse(LayoutDate, dstart)
	if err != nil {
		return "", err
	}

	repeatSlice := strings.Split(repeat, " ")
	switch repeatSlice[0] {
	case "d":
		if len(repeatSlice) != 2 {
			return "", fmt.Errorf("incorrect repeat format for days")
		}
		days, err := parseDays(repeatSlice[1])
		if err != nil {
			return "", err
		}
		return calculateNextDateByDays(date, now, days), nil
	case "y":
		if len(repeatSlice) != 1 {
			return "", fmt.Errorf("incorrect repeat format for years")
		}
		return calculateNextDateByYears(date, now), nil
	case "w":
		if len(repeatSlice) != 2 {
			return "", fmt.Errorf("incorrect repeat format for weeks")
		}
		weekDays, err := parseWeekDays(repeatSlice[1])
		if err != nil {
			return "", err
		}
		return calculateNextDateByWeekDays(date, now, weekDays), nil
	case "m":
		if len(repeatSlice) < 2 {
			return "", fmt.Errorf("incorrect repeat format for months")
		}
		days, months, err := parseMonthDays(repeatSlice[1:])
		if err != nil {
			return "", err
		}
		return calculateNextDateByMonthDays(date, now, days, months), nil
	default:
		return "", fmt.Errorf("incorrect character")
	}
}

func parseDays(daysStr string) (int, error) {
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 0, err
	}
	if days <= 0 || days > 364 {
		return 0, fmt.Errorf("incorrect day range")
	}
	return days, nil
}

func parseWeekDays(weekDay string) (map[int]bool, error) {
	weekdays := make(map[int]bool)
	for _, part := range strings.Split(weekDay, ",") {
		day, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("incorrect format day of the week")
		}
		if day < 1 || day > 7 {
			return nil, fmt.Errorf("the day of the week must be from 1 to 7")
		}
		weekdays[day] = true
	}
	return weekdays, nil
}

func parseMonthDays(monthDays []string) (map[int]bool, map[int]bool, error) {
	days := make(map[int]bool)
	months := make(map[int]bool)

	for _, part := range strings.Split(monthDays[0], ",") {
		day, err := strconv.Atoi(part)
		if err != nil {
			return nil, nil, fmt.Errorf("incorrect day of the month format")
		}
		if day < -2 || day == 0 || day > 31 {
			return nil, nil, fmt.Errorf("the day of the month must be from 1 to 31, -1 or -2")
		}
		days[day] = true
	}

	if len(monthDays) > 1 {
		for _, part := range strings.Split(monthDays[1], ",") {
			month, err := strconv.Atoi(part)
			if err != nil {
				return nil, nil, fmt.Errorf("incorrect month format")
			}
			if month < 1 || month > 12 {
				return nil, nil, fmt.Errorf("the month should be from 1 to 12")
			}
			months[month] = true
		}
	}

	return days, months, nil
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return y1 > y2 || (y1 == y2 && m1 > m2) || (y1 == y2 && m1 == m2 && d1 > d2)
}

func calculateNextDateByDays(start, now time.Time, days int) string {
	date := start
	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(LayoutDate)
}

func calculateNextDateByYears(start, now time.Time) string {
	date := start
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(LayoutDate)
}

func calculateNextDateByWeekDays(start, now time.Time, weekDays map[int]bool) string {
	date := start
	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			weekDay := int(date.Weekday())
			if weekDay == 0 {
				weekDay = 7
			}
			if weekDays[weekDay] {
				break
			}
		}
	}
	return date.Format(LayoutDate)
}

func calculateNextDateByMonthDays(start, now time.Time, days, months map[int]bool) string {
	date := start
	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			if len(months) > 0 {
				month := int(date.Month())
				if !months[month] {
					continue
				}
			}

			day := date.Day()
			lastDay := lastDayOfMonth(date)
			preLastDay := lastDay - 1

			if days[-1] && day == lastDay {
				break
			}
			if days[-2] && day == preLastDay {
				break
			}
			if days[day] {
				break
			}
		}
	}
	return date.Format(LayoutDate)
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

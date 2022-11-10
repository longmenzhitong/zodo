package zodo

import (
	"strconv"
	"strings"
	"time"
)

func CalcBetweenDays(t1 time.Time, t2 time.Time) (natureDays int, workDays int) {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	after := t1.After(t2)
	if after {
		tmp := t1
		t1 = t2
		t2 = tmp
	}

	for {
		if t1.Equal(t2) {
			break
		}

		natureDays++
		if t1.Weekday() != time.Saturday && t1.Weekday() != time.Sunday {
			workDays++
		}

		t1 = t1.Add(time.Hour * 24)
	}

	if after {
		natureDays *= -1
		workDays *= -1
	}
	return
}

func SimplifyTime(t string) string {
	if t == "" {
		return t
	}
	y := time.Now().Year()
	if strings.HasPrefix(t, strconv.Itoa(y)) {
		return t[len(strconv.Itoa(y))+1:]
	}
	return t
}

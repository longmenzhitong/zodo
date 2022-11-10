package zodo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func argsToStr(args []string) string {
	var str string
	for _, arg := range args {
		str += arg
	}
	return str
}

func argsToIds(args []string) (ids []int, err error) {
	ids = make([]int, 0)
	var id int
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		id, err = strconv.Atoi(arg)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

func argsToIdAndStr(args []string) (id int, str string, err error) {
	if len(args) != 2 {
		err = &InvalidInputError{
			Message: fmt.Sprintf("expect: [id] [str], got: %v", args),
		}
		return
	}
	id, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	str = strings.TrimSpace(args[1])
	return
}

func validateDeadline(ddl string) (string, error) {
	_, err := time.Parse(LayoutYearMonthDay, ddl)
	if err == nil {
		return ddl, nil
	}

	t, err := time.Parse(LayoutMonthDay, ddl)
	if err == nil {
		d := time.Date(time.Now().Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		return d.Format(LayoutYearMonthDay), nil
	}

	return "", &InvalidInputError{
		Message: fmt.Sprintf("deadline: %s", ddl),
	}
}

func validateRemind(rmd string) (string, error) {
	_, err := time.Parse(LayoutYearMonthDayHourMinute, rmd)
	if err == nil {
		return rmd, nil
	}

	now := time.Now()
	t, err := time.Parse(LayoutMonthDayHourMinute, rmd)
	if err == nil {
		d := time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		return d.Format(LayoutYearMonthDayHourMinute), nil
	}

	t, err = time.Parse(LayoutHourMinute, rmd)
	if err == nil {
		d := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		return d.Format(LayoutYearMonthDayHourMinute), nil
	}

	return "", &InvalidInputError{
		Message: fmt.Sprintf("remind: %s", rmd),
	}
}

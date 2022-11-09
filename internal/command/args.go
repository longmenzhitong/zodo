package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"zodo/internal/cst"
	"zodo/internal/errs"
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
		err = &errs.InvalidInputError{
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
	_, err := time.Parse(cst.LayoutYearMonthDay, ddl)
	if err == nil {
		return ddl, nil
	}

	t, err := time.Parse(cst.LayoutMonthDay, ddl)
	if err == nil {
		d := time.Date(time.Now().Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		return d.Format(cst.LayoutYearMonthDay), nil
	}

	return "", &errs.InvalidInputError{
		Message: fmt.Sprintf("deadline: %s", ddl),
	}
}

func validateRemind(rmd string) (string, error) {
	_, err := time.Parse(cst.LayoutYearMonthDayHourMinute, rmd)
	if err == nil {
		return rmd, nil
	}

	now := time.Now()
	t, err := time.Parse(cst.LayoutMonthDayHourMinute, rmd)
	if err == nil {
		d := time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		return d.Format(cst.LayoutYearMonthDayHourMinute), nil
	}

	t, err = time.Parse(cst.LayoutHourMinute, rmd)
	if err == nil {
		d := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		return d.Format(cst.LayoutYearMonthDayHourMinute), nil
	}

	return "", &errs.InvalidInputError{
		Message: fmt.Sprintf("remind: %s", rmd),
	}
}

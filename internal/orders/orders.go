package orders

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"sort"
	"strconv"
	"strings"
	"time"
	"zodo/internal/cst"
	"zodo/internal/errs"
	"zodo/internal/param"
	"zodo/internal/stdout"
	"zodo/internal/todos"
)

const (
	help          = "help"
	dailyReport   = "dr"
	add           = "add"
	_delete       = "del"
	modify        = "mod"
	rollback      = "rbk"
	transfer      = "trans"
	setDeadline   = "ddl"
	setRemark     = "rmk"
	setRemind     = "rmd"
	setLoopRemind = "rmd+"
	cancelRemind  = "rmd-"
	setChild      = "scd"
	addChild      = "acd"
	setPending    = "pend"
	setProcessing = "proc"
	setDone       = "done"
)

var orderDesc = map[string]string{
	help:          "view help info",
	dailyReport:   "send daily report email",
	add:           "add todo",
	_delete:       "delete todo",
	modify:        "modify todo",
	rollback:      "rollback to last version",
	transfer:      "transfer between file and redis",
	setDeadline:   "set deadline of todo",
	setRemark:     "set remark of todo",
	setRemind:     "set remind of todo",
	setLoopRemind: "set loop remind of todo",
	cancelRemind:  "cancel remind of todo",
	setChild:      "set child of todo",
	addChild:      "add child of todo",
	setPending:    "mark todo as pending",
	setProcessing: "mark todo as processing",
	setDone:       "mark todo as done",
}

func Handle(input string) error {
	input = strings.TrimSpace(input)
	order, val := parseInput(input)

	if order == help {
		orderList := make([]string, 0)
		for odr := range orderDesc {
			orderList = append(orderList, odr)
		}
		sort.Strings(orderList)
		rows := make([]table.Row, 0)
		for _, odr := range orderList {
			rows = append(rows, table.Row{odr, orderDesc[odr]})
		}
		stdout.PrintTable(table.Row{"Order", "Comment"}, rows)
		return nil
	}

	if order == dailyReport {
		return todos.DailyReport()
	}

	if order == add || param.ParentId != 0 || param.Deadline != "" {
		var content string
		if order == add {
			content = val
		} else {
			content = param.Input
		}
		id, err := todos.Add(content)
		if err != nil {
			return err
		}

		if param.ParentId != 0 {
			err = todos.SetChild(param.ParentId, []int{id}, true)
			if err != nil {
				return err
			}
		}

		if param.Deadline != "" {
			ddl, err := validateDeadline(param.Deadline)
			if err != nil {
				return err
			}
			todos.SetDeadline(id, ddl)
		}

		todos.Save()
		return nil
	}

	if order == _delete {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		todos.Delete(ids)
		todos.Save()
		return nil
	}

	if order == modify {
		id, content, err := parseIdAndStr(val)
		if err != nil {
			return err
		}
		todos.Modify(id, content)
		todos.Save()
		return nil
	}

	if order == rollback {
		todos.Rollback()
		return nil
	}

	if order == transfer {
		todos.Transfer()
		return nil
	}

	if order == setDeadline {
		id, ddl, err := parseIdAndStr(val)
		if err != nil {
			return err
		}
		ddl, err = validateDeadline(ddl)
		if err != nil {
			return err
		}
		todos.SetDeadline(id, ddl)
		todos.Save()
		return nil
	}

	if order == setRemark {
		id, remark, err := parseIdAndStr(val)
		if err != nil {
			return err
		}
		todos.SetRemark(id, remark)
		todos.Save()
		return nil
	}

	if order == setRemind || order == setLoopRemind {
		id, rmd, err := parseIdAndStr(val)
		if err != nil {
			return err
		}
		rmd, err = validateRemind(rmd)
		if err != nil {
			return err
		}
		err = todos.SetRemind(id, rmd, order == setLoopRemind)
		if err != nil {
			return err
		}
		todos.Save()
		return nil
	}

	if order == cancelRemind {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		todos.CancelRemind(ids)
		todos.Save()
		return nil
	}

	if order == setChild || order == addChild {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		if len(ids) < 2 {
			return &errs.InvalidInputError{
				Message: fmt.Sprintf("expect: %s [parentId] [childId], got: %s", order, input),
			}
		}
		err = todos.SetChild(ids[0], ids[1:], order == addChild)
		if err != nil {
			return err
		}
		todos.Save()
		return nil
	}

	if order == setPending {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todos.SetPending(id)
		}
		todos.Save()
		return nil
	}

	if order == setProcessing {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todos.SetProcessing(id)
		}
		todos.Save()
		return nil
	}

	if order == setDone {
		ids, err := parseIds(val)
		if err == nil {
			for _, id := range ids {
				todos.SetDone(id)
			}
			todos.Save()
			return nil
		}

		id, remark, err := parseIdAndStr(val)
		if err == nil {
			todos.SetDone(id)
			todos.SetRemark(id, remark)
			todos.Save()
			return nil
		}

		return err
	}

	// detail
	ids, err := parseIds(input)
	if err == nil && len(ids) > 0 {
		for _, id := range ids {
			todos.Detail(id)
			fmt.Println()
		}
		return nil
	}

	todos.List(input)
	return nil
}

func parseInput(input string) (order string, val string) {
	if input == "" {
		return
	}

	for odr := range orderDesc {
		if strings.HasPrefix(input, odr+" ") {
			order = odr
			val = strings.TrimSpace(strings.TrimPrefix(input, odr))
			return
		}
	}

	val = input
	return
}

func parseIds(val string) (ids []int, err error) {
	items := strings.Split(val, " ")
	ids = make([]int, 0)
	var id int
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		id, err = strconv.Atoi(item)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

func parseIdAndStr(val string) (id int, str string, err error) {
	i := strings.Index(val, " ")
	if i == -1 {
		err = &errs.InvalidInputError{
			Message: fmt.Sprintf("expect: [id] [str], got: %s", val),
		}
		return
	}
	id, err = strconv.Atoi(val[:i])
	if err != nil {
		return
	}
	str = strings.TrimSpace(val[i+1:])
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

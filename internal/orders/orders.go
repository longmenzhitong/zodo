package orders

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"zodo/internal/cst"
	"zodo/internal/errs"
	"zodo/internal/param"
	"zodo/internal/todos"
)

const (
	exit          = "exit"
	detail        = "cat"
	dailyReport   = "dr"
	add           = "add"
	_delete       = "del"
	modify        = "mod"
	transfer      = "trans"
	setDeadline   = "ddl"
	setRemark     = "rmk"
	setChild      = "scd"
	addChild      = "acd"
	setPending    = "pend"
	setProcessing = "proc"
	setDone       = "done"
)

var allOrders = []string{
	exit, detail, dailyReport, add, _delete, modify, transfer,
	setDeadline, setRemark, setChild, addChild, setPending, setProcessing, setDone,
}

func Handle(input string) error {
	input = strings.TrimSpace(input)
	order, val := parseInput(input)

	if order == exit {
		fmt.Println("Bye.")
		os.Exit(0)
	}

	if order == detail {
		id, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		todos.Detail(id)
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

	if order == _delete || param.Delete {
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

	if order == transfer {
		todos.Transfer()
		return nil
	}

	if order == setDeadline {
		id, deadline, err := parseDeadline(val)
		if err != nil {
			return err
		}
		todos.SetDeadline(id, deadline)
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

	if order == setChild || order == addChild {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		if len(ids) < 2 {
			return &errs.InvalidInputError{
				Input:   input,
				Message: fmt.Sprintf("expect: %s [parentId] [childId]", order),
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
		if err != nil {
			return err
		}
		for _, id := range ids {
			todos.SetDone(id)
		}
		todos.Save()
		return nil
	}

	id, err := strconv.Atoi(input)
	if err == nil {
		todos.Detail(id)
		return nil
	}

	todos.List(input)
	return nil
}

func parseInput(input string) (order string, val string) {
	if input == "" {
		return
	}

	for _, odr := range allOrders {
		if strings.HasPrefix(input, odr) {
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
			Input:   val,
			Message: "expect: [id] [str]",
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

func parseDeadline(val string) (id int, ddl string, err error) {
	id, ddl, err = parseIdAndStr(val)
	if err != nil {
		return
	}

	ddl, err = validateDeadline(ddl)
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
		Input:   "deadline",
		Message: ddl,
	}
}

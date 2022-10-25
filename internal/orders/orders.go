package orders

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"zodo/internal/backup"
	"zodo/internal/cst"
	"zodo/internal/errs"
	"zodo/internal/todo"
)

const (
	exit = "exit"
	pull = "pull"
	push = "push"
)

const (
	// 查询指令
	list   = "ll"
	detail = "cat"
)

const (
	// 修改值的指令
	add         = "add"
	modify      = "mod"
	setDeadline = "ddl"
	setRemark   = "rmk"
)

const (
	// 修改状态的指令
	setPending    = "pend"
	setProcessing = "proc"
	setDone       = "done"
	setDeleted    = "del"
)

func Handle(input string) error {
	input = strings.TrimSpace(input)
	order, val := parseInput(input)

	if order == exit {
		fmt.Println("Bye.")
		os.Exit(0)
	}

	if order == pull {
		return backup.Pull()
	}

	if order == push {
		return backup.Push()
	}

	if order == list {
		todo.List()
		return nil
	}

	if order == detail {
		id, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		todo.Detail(id)
		return nil
	}

	if order == add {
		todo.Add(val)
		return nil
	}

	if order == modify {
		id, content, err := parseIdAndStr(val)
		if err != nil {
			return err
		}
		todo.Modify(id, content)
		return nil
	}

	if order == setDeadline {
		id, deadline, err := parseDeadline(val)
		if err != nil {
			return err
		}
		todo.Deadline(id, deadline)
		return nil
	}

	if order == setRemark {
		id, remark, err := parseIdAndStr(val)
		if err != nil {
			return err
		}
		todo.Remark(id, remark)
		return nil
	}

	if order == setPending {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Pending(id)
		}
		return nil
	}

	if order == setProcessing {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Processing(id)
		}
		return nil
	}

	if order == setDone {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Done(id)
		}
		return nil
	}

	if order == setDeleted {
		ids, err := parseIds(val)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Delete(id)
		}
		return nil
	}

	// todo help
	// todo hint

	id, err := strconv.Atoi(input)
	if err == nil {
		todo.Detail(id)
		return nil
	}

	todo.List()
	return nil
}

func parseInput(input string) (order string, val string) {
	if input == "" {
		return
	}
	i := strings.Index(input, " ")
	if i == -1 {
		order = input
	} else {
		order = input[:i]
		val = input[i+1:]
	}
	return strings.TrimSpace(order), strings.TrimSpace(val)
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
		err = &errs.InvalidInputError{Input: val}
		return
	}
	id, err = strconv.Atoi(val[:i])
	if err != nil {
		return
	}
	str = strings.TrimSpace(val[i+1:])
	return
}

func parseDeadline(val string) (id int, deadline string, err error) {
	id, deadline, err = parseIdAndStr(val)
	if err != nil {
		return
	}

	_, err = time.Parse(cst.LayoutYearMonthDay, deadline)
	if err == nil {
		return
	}

	t, err := time.Parse(cst.LayoutMonthDay, deadline)
	if err == nil {
		d := time.Date(time.Now().Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		deadline = d.Format(cst.LayoutYearMonthDay)
		return
	}

	return
}

package orders

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"zodo/internal/cst"
	"zodo/internal/errs"
	"zodo/internal/todo"
)

const (
	exit = "exit"
)

const (
	// 查询指令
	list         = "ll"
	prefixDetail = "cat "
)

const (
	// 修改值的指令
	prefixAdd      = "add "
	prefixModify   = "mod "
	prefixDeadline = "ddl "
	prefixRemark   = "rmk "
)

const (
	// 修改状态的指令
	prefixPending    = "pend "
	prefixProcessing = "proc "
	prefixDone       = "done "
	prefixDelete     = "rm "
)

func Handle(input string) error {
	defer todo.Save()

	if strings.TrimSpace(input) == exit {
		todo.Save()
		fmt.Println("Bye.")
		os.Exit(0)
	}

	if strings.TrimSpace(input) == list {
		todo.List()
		return nil
	}

	if strings.HasPrefix(input, prefixDetail) {
		id, err := parseId(input, prefixDetail)
		if err != nil {
			return err
		}
		todo.Detail(id)
		return nil
	}

	if strings.HasPrefix(input, prefixAdd) {
		content, err := parseStr(input, prefixAdd)
		if err != nil {
			return err
		}
		todo.Add(content)
		return nil
	}

	if strings.HasPrefix(input, prefixModify) {
		id, content, err := parseIdAndStr(input, prefixModify)
		if err != nil {
			return err
		}
		todo.Modify(id, content)
		return nil
	}

	if strings.HasPrefix(input, prefixDeadline) {
		id, deadline, err := parseDeadline(input)
		if err != nil {
			return err
		}
		todo.Deadline(id, deadline)
		return nil
	}

	if strings.HasPrefix(input, prefixRemark) {
		id, remark, err := parseIdAndStr(input, prefixRemark)
		if err != nil {
			return err
		}
		todo.Remark(id, remark)
		return nil
	}

	if strings.HasPrefix(input, prefixPending) {
		id, err := parseId(input, prefixPending)
		if err != nil {
			return err
		}
		todo.Pending(id)
		return nil
	}

	if strings.HasPrefix(input, prefixProcessing) {
		id, err := parseId(input, prefixProcessing)
		if err != nil {
			return err
		}
		todo.Processing(id)
		return nil
	}

	if strings.HasPrefix(input, prefixDone) {
		id, err := parseId(input, prefixDone)
		if err != nil {
			return err
		}
		todo.Done(id)
		return nil
	}

	if strings.HasPrefix(input, prefixDelete) {
		id, err := parseId(input, prefixDelete)
		if err != nil {
			return err
		}
		todo.Delete(id)
		return nil
	}

	// todo help
	// todo detail
	// todo hint

	todo.List()
	return nil
}

func parseDeadline(input string) (id int, deadline string, err error) {
	id, deadline, err = parseIdAndStr(input, prefixDeadline)
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

func parseId(input, prefix string) (id int, err error) {
	order := strings.TrimSpace(strings.TrimPrefix(input, prefix))
	return strconv.Atoi(order)
}

func parseStr(input, prefix string) (str string, err error) {
	str = strings.TrimSpace(strings.TrimPrefix(input, prefix))
	if str == "" {
		err = &errs.InvalidInputError{Input: ""}
	}
	return
}

func parseIdAndStr(input, prefix string) (id int, str string, err error) {
	order := strings.TrimSpace(strings.TrimPrefix(input, prefix))
	i := strings.Index(order, " ")
	if i == -1 {
		err = &errs.InvalidInputError{Input: input}
		return
	}
	id, err = strconv.Atoi(order[:i])
	if err != nil {
		return
	}
	str = strings.TrimSpace(order[i+1:])
	if str == "" {
		err = &errs.InvalidInputError{Input: input}
	}
	return
}

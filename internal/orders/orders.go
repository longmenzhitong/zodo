package orders

import (
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
	list = "ll"
)

const (
	prefixAdd      = "add "
	prefixModify   = "mod "
	prefixDeadline = "ddl "
	prefixDone     = "done "
	prefixPending  = "hang "
	prefixDelete   = "del "
)

func Handle(input string) error {
	defer todo.Save()

	if strings.TrimSpace(input) == exit {
		todo.Save()
		os.Exit(0)
	}

	if strings.TrimSpace(input) == list {
		todo.List()
		return nil
	}

	if strings.HasPrefix(input, prefixAdd) {
		content, err := ParseAdd(input)
		if err != nil {
			return err
		}
		todo.Add(content)
		return nil
	}

	if strings.HasPrefix(input, prefixModify) {
		id, content, err := ParseModify(input)
		if err != nil {
			return err
		}
		todo.Modify(id, content)
		return nil
	}

	if strings.HasPrefix(input, prefixDeadline) {
		id, deadline, err := ParseDeadline(input)
		if err != nil {
			return err
		}
		todo.Deadline(id, deadline)
		return nil
	}

	if strings.HasPrefix(input, prefixPending) {
		id, err := ParsePending(input)
		if err != nil {
			return err
		}
		todo.Pending(id)
		return nil
	}

	if strings.HasPrefix(input, prefixDone) {
		id, err := ParseDone(input)
		if err != nil {
			return err
		}
		todo.Done(id)
		return nil
	}

	if strings.HasPrefix(input, prefixDelete) {
		id, err := ParseDelete(input)
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

func ParseAdd(input string) (content string, err error) {
	content = strings.TrimSpace(strings.TrimPrefix(input, prefixAdd))
	if content == "" {
		err = &errs.InvalidInputError{Input: ""}
	}
	return
}

func ParseModify(input string) (id int, content string, err error) {
	return parseIdAndStr(input, prefixModify)
}

func ParseDeadline(input string) (id int, deadline string, err error) {
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

func ParsePending(input string) (id int, err error) {
	return parseId(input, prefixPending)
}

func ParseDone(input string) (id int, err error) {
	return parseId(input, prefixDone)
}

func ParseDelete(input string) (id int, err error) {
	return parseId(input, prefixDelete)
}

func parseId(input, prefix string) (id int, err error) {
	order := strings.TrimSpace(strings.TrimPrefix(input, prefix))
	return strconv.Atoi(order)
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

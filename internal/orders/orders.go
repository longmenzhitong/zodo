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
	prefixDelete     = "del "
)

func Handle(input string) error {
	defer todo.Save()

	if strings.TrimSpace(input) == exit {
		todo.Save()
		fmt.Println("Bye.")
		os.Exit(0)
	}

	if strings.TrimSpace(input) == pull {
		return backup.Pull()
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
		todo.Add(parseStr(input, prefixAdd))
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
		ids, err := parseIds(input, prefixPending)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Pending(id)
		}
		return nil
	}

	if strings.HasPrefix(input, prefixProcessing) {
		ids, err := parseIds(input, prefixProcessing)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Processing(id)
		}
		return nil
	}

	if strings.HasPrefix(input, prefixDone) {
		ids, err := parseIds(input, prefixDone)
		if err != nil {
			return err
		}
		for _, id := range ids {
			todo.Done(id)
		}
		return nil
	}

	if strings.HasPrefix(input, prefixDelete) {
		ids, err := parseIds(input, prefixDelete)
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

	id, err := parseId(input, "")
	if err == nil {
		todo.Detail(id)
		return nil
	}

	todo.List()
	return nil
}

func parseId(input, prefix string) (id int, err error) {
	return strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(input, prefix)))
}

func parseIds(input, prefix string) (ids []int, err error) {
	order := strings.TrimPrefix(input, prefix)
	items := strings.Split(order, " ")
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

func parseStr(input, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(input, prefix))
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
	return
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

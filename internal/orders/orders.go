package orders

import (
	"strconv"
	"strings"
	"time"
	"zodo/internal/cst"
	"zodo/internal/errs"
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
	prefixAbandon  = "abd "
	prefixDelete   = "del "
)

func IsExit(input string) bool {
	return strings.TrimSpace(input) == exit
}

func IsList(input string) bool {
	return strings.TrimSpace(input) == list
}

func IsAdd(input string) bool {
	return strings.HasPrefix(input, prefixAdd)
}

func ParseAdd(input string) (content string, err error) {
	content = strings.TrimSpace(strings.TrimPrefix(input, prefixAdd))
	if content == "" {
		err = &errs.InvalidInputError{Input: ""}
	}
	return
}

func IsModify(input string) bool {
	return strings.HasPrefix(input, prefixModify)
}

func ParseModify(input string) (id int, content string, err error) {
	return parseIdAndStr(input, prefixModify)
}

func IsDeadline(input string) bool {
	return strings.HasPrefix(input, prefixDeadline)
}

func ParseDeadline(input string) (id int, deadline string, err error) {
	id, deadline, err = parseIdAndStr(input, prefixDeadline)
	if err != nil {
		return
	}
	_, err = time.Parse(cst.LayoutDate, deadline)
	return
}

func IsPending(input string) bool {
	return strings.HasPrefix(input, prefixPending)
}

func ParsePending(input string) (id int, err error) {
	return parseSingleId(input, prefixPending)
}

func IsDone(input string) bool {
	return strings.HasPrefix(input, prefixDone)
}

func ParseDone(input string) (id int, err error) {
	return parseSingleId(input, prefixDone)
}

func IsAbandon(input string) bool {
	return strings.HasPrefix(input, prefixAbandon)
}

func ParseAbandon(input string) (id int, err error) {
	return parseSingleId(input, prefixAbandon)
}

func IsDelete(input string) bool {
	return strings.HasPrefix(input, prefixDelete)
}

func ParseDelete(input string) (id int, err error) {
	return parseSingleId(input, prefixDelete)
}

func parseSingleId(input, prefix string) (id int, err error) {
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

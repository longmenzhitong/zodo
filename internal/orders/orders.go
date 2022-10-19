package orders

import (
	"strconv"
	"strings"
	"zodo/internal/errs"
)

const (
	exit = "exit"
	list = "ll"
)

const (
	prefixAdd     = "add "
	prefixModify  = "mod "
	prefixDone    = "done "
	prefixPending = "pending "
	prefixAbandon = "abandon "
	prefixDelete  = "del "
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
	order := strings.TrimSpace(strings.TrimPrefix(input, prefixModify))
	i := strings.Index(order, " ")
	if i == -1 {
		err = &errs.InvalidInputError{Input: input}
		return
	}
	id, err = strconv.Atoi(order[:i])
	if err != nil {
		return
	}
	content = strings.TrimSpace(order[i+1:])
	if content == "" {
		err = &errs.InvalidInputError{Input: input}
	}
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

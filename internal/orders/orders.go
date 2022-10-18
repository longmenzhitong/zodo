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

func IsExit(input string) error {
	if input != exit {
		return &errs.WrongOrderError{}
	}
	return nil
}

func IsList(input string) error {
	if input != list {
		return &errs.WrongOrderError{}
	}
	return nil
}

func IsAdd(input string) (string, error) {
	if !strings.HasPrefix(input, prefixAdd) {
		return "", &errs.WrongOrderError{}
	}
	content := strings.TrimSpace(strings.TrimPrefix(input, prefixAdd))
	if content == "" {
		return "", &errs.InvalidInputError{Input: ""}
	}
	return content, nil
}

func IsModify(input string) (int, string, error) {
	if !strings.HasPrefix(input, prefixModify) {
		return 0, "", &errs.WrongOrderError{}
	}
	order := strings.TrimSpace(strings.TrimPrefix(input, prefixModify))
	items := strings.Split(order, " ")
	if len(items) != 2 {
		return 0, "", &errs.InvalidInputError{Input: input}
	}
	id, err := strconv.Atoi(items[0])
	if err != nil {
		return 0, "", err
	}
	content := items[1]
	if content == "" {
		return 0, "", &errs.InvalidInputError{Input: input}
	}
	return id, content, nil
}

func IsPending(input string) (int, error) {
	return isModifyStatus(input, prefixPending)
}

func IsDone(input string) (int, error) {
	return isModifyStatus(input, prefixDone)
}

func IsAbandon(input string) (int, error) {
	return isModifyStatus(input, prefixAbandon)
}

func IsDelete(input string) (int, error) {
	return isModifyStatus(input, prefixDelete)
}

func isModifyStatus(input, prefix string) (int, error) {
	if !strings.HasPrefix(input, prefix) {
		return 0, &errs.WrongOrderError{}
	}
	order := strings.TrimSpace(strings.TrimPrefix(input, prefix))
	id, err := strconv.Atoi(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

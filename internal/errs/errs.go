package errs

import "fmt"

type InvalidInputError struct {
	Input string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Input)
}

type WrongOrderError struct {
}

func (e *WrongOrderError) Error() string {
	return "wrong order"
}

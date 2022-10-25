package errs

import "fmt"

type InvalidInputError struct {
	Input   string
	Message string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input [%s]: %s, ", e.Input, e.Message)
}

type NotFoundError struct {
	Target  string
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("[%s] not found: %s", e.Target, e.Message)
}

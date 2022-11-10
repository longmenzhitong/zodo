package zodo

import "fmt"

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s, ", e.Message)
}

type InvalidConfigError struct {
	Message string
}

func (e *InvalidConfigError) Error() string {
	return fmt.Sprintf("invalid config: %s, ", e.Message)
}

type NotFoundError struct {
	Target  string
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Target, e.Message)
}

type CancelledError struct {
}

func (e *CancelledError) Error() string {
	return "Cancelled."
}

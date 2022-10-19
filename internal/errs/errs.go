package errs

import "fmt"

type InvalidInputError struct {
	Input string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Input)
}

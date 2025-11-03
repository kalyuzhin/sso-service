package error

import "fmt"

// NestedError – ...
type NestedError struct {
	Message string
	Err     error
}

// Error – ...
func (e *NestedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}

	return e.Message
}

// WrapErr – ...
func WrapErr(err error, msg string) *NestedError {
	return &NestedError{
		Message: msg,
		Err:     err,
	}
}

// New – ...
func New(msg string) *NestedError {
	return &NestedError{Message: msg}
}

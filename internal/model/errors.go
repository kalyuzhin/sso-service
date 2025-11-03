package model

import "fmt"

// NestedError – ...
type NestedError struct {
	Message string
	Err     error
}

// Error – ...
func (e *NestedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
}

// WrapErr – ...
func WrapErr(err error, msg string) *NestedError {
	return &NestedError{
		Message: msg,
		Err:     err,
	}
}

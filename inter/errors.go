package inter

import "fmt"

// SpecError is the error type returned if there are problems with a
// specification
type SpecError struct {
	Message string
}

// SpecErrorf creates a SpecError from a format string and arguments
func SpecErrorf(format string, args ...interface{}) *SpecError {
	return &SpecError{fmt.Sprintf(format, args...)}
}

// Error returns a string representation of a SpecError, implementing
// the error interface.
func (err *SpecError) Error() string {
	return err.Message
}

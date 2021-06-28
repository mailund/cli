package params

import "fmt"

// ParseError is the type that the parser will return on errors.
// It implements the error interface.
type ParseError struct {
	Message string
}

// ParseErrorf creates a ParseError from a format string and arguments.
func ParseErrorf(format string, args ...interface{}) *ParseError {
	return &ParseError{fmt.Sprintf(format, args...)}
}

// Error returns a string from a ParseError, implementing the error
// interface.
func (err *ParseError) Error() string {
	return err.Message
}

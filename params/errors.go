package params

import "fmt"

type ParseError struct {
	Message string
}

func ParseErrorf(format string, args ...interface{}) *ParseError {
	return &ParseError{fmt.Sprintf(format, args...)}
}

func (err *ParseError) Error() string {
	return err.Message
}

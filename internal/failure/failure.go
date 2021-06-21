package failure

import "os"

// This is only used for testing. It allows me
// to mock terminating the program

// DefaultFailure is to exit with value 2
var DefaultFailure = func() { os.Exit(2) }

// Failure, if overwritten, will do something else
var Failure = DefaultFailure

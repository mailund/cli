package failure

import "os"

// This is only used for testing. It allows me
// to mock terminating the program

var DefaultFailure = func() { os.Exit(2) }
var Failure = DefaultFailure

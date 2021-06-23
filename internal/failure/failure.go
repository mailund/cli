package failure

import "os"

// This is only used for testing. It allows me
// to mock terminating the program

const exitStatus = 2

// DefaultFailure is to exit with value 2
var DefaultFailure = func() { os.Exit(exitStatus) }

// Failure is called when the package wants the program to terminate.
// You can overwrite it when mocking failure in tests.
var Failure = DefaultFailure

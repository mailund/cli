package failure

type (
	/* ErrorHandling flags for determining how the argument
	parser should handle errors.

		- ExitOnError: Terminate the program on error.
		- ContinueOnError: Report the error as an error object.
		- PanicOnError: Panics if there is an error.
	*/
	ErrorHandling int
)

const (
	ContinueOnError ErrorHandling = iota // ContinueOnError means that parsing will return an error
	ExitOnError                   = iota // ExitOnError means that parsing will exit the program
	PanicOnError                  = iota // PanicOnError means we raise a panic on errors
)

// Implementations of Values interfaces for wrapping flags and positional
// arguments
package vals

type PosValue interface {
	Set(string) error
}

type FlagValue interface {
	String() string
	Set(string) error
}

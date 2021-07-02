// Package interfaces implements public interfaces for functionality buried
// deeper in the internal packages
package interfaces

// PosValue is the interface that positional arguments must implement
type PosValue interface {
	Set(string) error // Should set the value from a string
}

// VariadicValue is the interface for variadic positional arguments
type VariadicValue interface {
	Set([]string) error // Should set the value from a slice of strings
}

// FlagValue is the interface that flag arguments must implement
type FlagValue interface {
	String() string   // Should return a string representation of the value
	Set(string) error // Should set the value from a string
}

// NoValueFlag is used to indicate that a flag doesn't take any values, and
// that it is an error to provide one. Their Set() method will be called with
// the empty string instead.
type NoValueFlag interface {
	NoValueFlag() bool
}

// DefaultValueFlag is used to indicate that a flag doesn't need a value, but can take one.
// This differs from NoValueFlag where it is an error to provide one. With DefaultValueFlag,
// the Default() function should return the string that Set() is called with if no value
// is provided
type DefaultValueFlag interface {
	DefaultValueFlag() string
}

// ArgumentDescription provides a value a way to add to the description string for a flag or positional.
type ArgumentDescription interface {
	ArgumentDescription(flag bool, descr string) string // Modify or add to the description string
}

// FlagValueDescription provides a value a way to add to the "value" string of a flag.
type FlagValueDescription interface {
	FlagValueDescription() string // Modify or add to the description string
}

// Validator is an interface that is run after parameters are initialised
// but before they are parsed, and can be used to check consistency
// of default values. The flag is true if validating a flag and false otherwise.
type Validator interface {
	Validate(flag bool) error // Should return nil if everything is fine, or an error otherwise
}

// Prepare is used after parsing and before running a command.
type Prepare interface {
	PrepareValue() error // Called after parsing and before we run a command
}

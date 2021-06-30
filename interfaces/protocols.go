// interfaces implements public interfaces for functionality buried
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

// BoolFlag is used to indicate that a flag doesn't need a value, but can take one.
// This differs from NoValueFlag where it is an error to provide one. With DefaultValueFlag,
// the Default() function should return the string that Set() is called with if no value
// is provided
type DefaultValueFlag interface {
	Default() string
}

// Validate is an interface that is run after parameters are initialised
// but before they are parsed, and can be used to check consistency
// of default values
type Validator interface {
	Validate() error // Should return nil if everything is fine, or an error otherwise
}

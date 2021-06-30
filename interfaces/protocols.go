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

// BoolFlag is used to indicate if a flag needs an argument. For boolean flags,
// `--flag` can stand alone, but otherwise flags must have an argument as `--flag=arg`
// or `--flag arg`
type BoolFlag interface {
	IsBoolFlag() bool // Should return true if the flag doesn't need an argument
}

// Validate is an interface that is run after parameters are initialised
// but before they are parsed, and can be used to check consistency
// of default values
type Validator interface {
	Validate() error // Should return nil if everything is fine, or an error otherwise
}

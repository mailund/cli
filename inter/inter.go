// inter implements public interfaces for functionality buried
// deeper in the internal packages
package inter

// PosValue is the interface that positional arguments must implement
type PosValue interface {
	Set(string) error
}

// FlagValue is the interface that flag arguments must implement
type FlagValue interface {
	String() string
	Set(string) error
}

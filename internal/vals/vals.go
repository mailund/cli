// Package vals contain interfaces for errors and for
// command line value protocols.
package vals

import (
	"reflect"

	"github.com/mailund/cli/interfaces"
)

//go:generate go run gen/genvals.go genvals

// DefaultValueFlag gives boolean values a truth default if used as
// a flag with no value
func (val *BoolValue) DefaultValueFlag() string { return "true" }

type valConstructor func(reflect.Value) interfaces.FlagValue

var valsConstructors = map[reflect.Type]valConstructor{}

// AsFlagValue attempts to turn a value into a FlagValue interface
func AsFlagValue(val reflect.Value) interfaces.FlagValue {
	if cast, ok := val.Interface().(interfaces.FlagValue); ok {
		return cast
	}

	if cons, ok := valsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

// AsPosValue attempts to turn a value into a PosValue interface
func AsPosValue(val reflect.Value) interfaces.PosValue {
	if cast, ok := val.Interface().(interfaces.PosValue); ok {
		return cast
	}

	if cons, ok := valsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

type varValConstructor func(reflect.Value) interfaces.VariadicValue

var varValsConstructors = map[reflect.Type]varValConstructor{}

// AsVariadicValue attempts to turn a value into a variadic value
// interface
func AsVariadicValue(val reflect.Value) interfaces.VariadicValue {
	if cast, ok := val.Interface().(interfaces.VariadicValue); ok {
		return cast
	}

	if cons, ok := varValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

// Implementations of Values interfaces for wrapping flags and positional
// arguments
package vals

import (
	"reflect"

	"github.com/mailund/cli/interfaces"
)

//go:generate go run gen/genvals.go genvals

// We explicitly write this method so boolean flags do not need arguments
func (val *BoolValue) Default() string { return "true" }

type ValConstructor func(reflect.Value) interfaces.FlagValue

var ValsConstructors = map[reflect.Type]ValConstructor{}

func AsFlagValue(val reflect.Value) interfaces.FlagValue {
	if cast, ok := val.Interface().(interfaces.FlagValue); ok {
		return cast
	}

	if cons, ok := ValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

func AsPosValue(val reflect.Value) interfaces.PosValue {
	if cast, ok := val.Interface().(interfaces.PosValue); ok {
		return cast
	}

	if cons, ok := ValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

type VarValConstructor func(reflect.Value) interfaces.VariadicValue

var VarValsConstructors = map[reflect.Type]VarValConstructor{}

func AsVariadicValue(val reflect.Value) interfaces.VariadicValue {
	if cast, ok := val.Interface().(interfaces.VariadicValue); ok {
		return cast
	}

	if cons, ok := VarValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

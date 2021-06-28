// Implementations of Values interfaces for wrapping flags and positional
// arguments
package vals

import (
	"reflect"

	"github.com/mailund/cli/inter"
)

//go:generate sh -c "go run gen/genvals.go > genvals.go"

type ValConstructor func(reflect.Value) inter.FlagValue

var ValsConstructors = map[reflect.Type]ValConstructor{}

func AsValue(val reflect.Value) inter.FlagValue {
	if cast, ok := val.Interface().(inter.FlagValue); ok {
		return cast
	}

	if cons, ok := ValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

type VarValConstructor func(reflect.Value) inter.VariadicValue

var VarValsConstructors = map[reflect.Type]VarValConstructor{}

func AsVariadicValue(val reflect.Value) inter.VariadicValue {
	if cast, ok := val.Interface().(inter.VariadicValue); ok {
		return cast
	}

	if cons, ok := VarValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

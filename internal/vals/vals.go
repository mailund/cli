// Implementations of Values interfaces for wrapping flags and positional
// arguments
package vals

import (
	"fmt"
	"reflect"

	"github.com/mailund/cli/inter"
)

//go:generate sh -c "go run gen/genvals.go > genvals.go"

type ValConstructor func(reflect.Value) inter.FlagValue

var ValsConstructors = map[reflect.Type]ValConstructor{}

func AsValue(val reflect.Value) inter.FlagValue {
	if cast, ok := val.Interface().(inter.FlagValue); ok {
		fmt.Printf("[%q, %q] satisfy the interface\n", val, val.Type())
		return cast
	}

	if cons, ok := ValsConstructors[val.Type()]; ok {
		return cons(val)
	}

	return nil
}

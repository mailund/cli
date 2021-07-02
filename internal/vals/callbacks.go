package vals

import (
	"reflect"

	"github.com/mailund/cli/interfaces"
)

// Function callbacks are special values; we don't need to generate them
// since it is easier just to write the little code there is in the first place
type (
	// FuncNoValue wraps functions that do not take parameter arguments
	FuncNoValue func() error
	// FuncValue wraps functions with a single string argument
	FuncValue func(string) error
	// VariadicFuncValue wraps functions with string slice argument
	VariadicFuncValue func([]string) error
)

// Set implements the PosValue/FlagValue interface
func (f FuncNoValue) Set(_ string) error { return f() }

// String implements the FlagValue interface
func (f FuncNoValue) String() string { return "" }

// NoValueFlag implements the no values interface
func (f FuncNoValue) NoValueFlag() bool { return true }

// Validate implements the validator interface and checks that callbacks are not nil.
func (f FuncNoValue) Validate(bool) error { // FIXME: the validators are the same, but need generics to join them
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

// Set implements the PosValue/FlagValue interface
func (f FuncValue) Set(x string) error { return f(x) }

// String implements the FlagValue interface
func (f FuncValue) String() string { return "" }

// Validate implements the validator interface and checks that callbacks are not nil.
func (f FuncValue) Validate(bool) error {
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

// Set implements the PosValue/FlagValue interface
func (f VariadicFuncValue) Set(x []string) error { return f(x) }

// Validate implements the validator interface and checks that callbacks are not nil.
func (f VariadicFuncValue) Validate(bool) error {
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

type (
	ncb   = func()
	ncbi  = func(interface{})
	ncbe  = func() error
	ncbei = func(interface{}) error

	cb  = func(string) error
	cbi = func(string, interface{}) error

	vcb  = func([]string) error
	vcbi = func([]string, interface{}) error

	callbackWrapper         func(*reflect.Value, interface{}) interfaces.FlagValue
	variadicCallbackWrapper func(*reflect.Value, interface{}) interfaces.VariadicValue
)

func ncbConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(ncb) // If we call this func, we know we have the right type
	return FuncNoValue(func() error { f(); return nil })
}

func ncbiConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(ncbi) // If we call this func, we know we have the right type
	return FuncNoValue(func() error { f(i); return nil })
}

func ncbeConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(ncbe) // If we call this func, we know we have the right type
	return FuncNoValue(f)
}

func ncbeiConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(ncbei) // If we call this func, we know we have the right type
	return FuncNoValue(func() error { return f(i) })
}

func cbConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(cb) // If we call this func, we know we have the right type
	return FuncValue(f)
}

func cbiConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(cbi) // If we call this func, we know we have the right type
	return FuncValue(func(x string) error { return f(x, i) })
}

func vcbConstructor(val *reflect.Value, i interface{}) interfaces.VariadicValue {
	f, _ := val.Interface().(vcb) // If we call this func, we know we have the right type
	return VariadicFuncValue(f)
}

func vcbiConstructor(val *reflect.Value, i interface{}) interfaces.VariadicValue {
	f, _ := val.Interface().(vcbi) // If we call this func, we know we have the right type
	return VariadicFuncValue(func(x []string) error { return f(x, i) })
}

var callbackWrappers = map[reflect.Type]callbackWrapper{
	reflect.TypeOf(ncb(nil)):   ncbConstructor,
	reflect.TypeOf(ncbi(nil)):  ncbiConstructor,
	reflect.TypeOf(ncbe(nil)):  ncbeConstructor,
	reflect.TypeOf(ncbei(nil)): ncbeiConstructor,
	reflect.TypeOf(cb(nil)):    cbConstructor,
	reflect.TypeOf(cbi(nil)):   cbiConstructor,
}

var variadicCallbackWrappers = map[reflect.Type]variadicCallbackWrapper{
	reflect.TypeOf(vcb(nil)):  vcbConstructor,
	reflect.TypeOf(vcbi(nil)): vcbiConstructor,
}

// AsCallback returns a value as a wrapped callback function
func AsCallback(val *reflect.Value, i interface{}) interfaces.FlagValue {
	if wrap, ok := callbackWrappers[val.Type()]; ok {
		return wrap(val, i)
	}

	return nil // couldn't convert
}

// AsVariadicCallback returns a value as a wrapped variadic callback function
func AsVariadicCallback(val *reflect.Value, i interface{}) interfaces.VariadicValue {
	if wrap, ok := variadicCallbackWrappers[val.Type()]; ok {
		return wrap(val, i)
	}

	return nil // couldn't convert
}

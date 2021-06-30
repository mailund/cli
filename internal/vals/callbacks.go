package vals

import (
	"reflect"

	"github.com/mailund/cli/interfaces"
)

// Function callbacks are special values; we don't need to generate them
// since it is easier just to write the little code there is in the first place
type (
	FuncBoolValue     func() error         // Wrapping functions that do not take parameter arguments
	FuncValue         func(string) error   // Wrapping functions with a single string argument
	VariadicFuncValue func([]string) error // Wrapping functions with string slice argument
)

func (f FuncBoolValue) Set(_ string) error { return f() }
func (f FuncBoolValue) String() string     { return "" }
func (f FuncBoolValue) IsBoolFlag() bool   { return true }
func (f FuncBoolValue) Validate() error { // FIXME: the validators are the same, but need generics to join them
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

func (f FuncValue) Set(x string) error { return f(x) }
func (f FuncValue) String() string     { return "" }
func (f FuncValue) Validate() error {
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

func (f VariadicFuncValue) Set(x []string) error { return f(x) }
func (f VariadicFuncValue) Validate() error {
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

type (
	BCB   = func()
	BCBI  = func(interface{})
	BCBE  = func() error
	BCBIE = func(interface{}) error

	CB  = func(string) error
	CBI = func(string, interface{}) error

	VCB  = func([]string) error
	VCBI = func([]string, interface{}) error

	CallbackWrapper         func(*reflect.Value, interface{}) interfaces.FlagValue
	VariadicCallbackWrapper func(*reflect.Value, interface{}) interfaces.VariadicValue
)

func BCBConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(BCB) // If we call this func, we know we have the right type
	return FuncBoolValue(func() error { f(); return nil })
}

func BCBIConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(BCBI) // If we call this func, we know we have the right type
	return FuncBoolValue(func() error { f(i); return nil })
}

func BCBEConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(BCBE) // If we call this func, we know we have the right type
	return FuncBoolValue(f)
}

func BCBIEConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(BCBIE) // If we call this func, we know we have the right type
	return FuncBoolValue(func() error { return f(i) })
}

func CBConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(CB) // If we call this func, we know we have the right type
	return FuncValue(f)
}

func CBIConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(CBI) // If we call this func, we know we have the right type
	return FuncValue(func(x string) error { return f(x, i) })
}

func VCBConstructor(val *reflect.Value, i interface{}) interfaces.VariadicValue {
	f, _ := val.Interface().(VCB) // If we call this func, we know we have the right type
	return VariadicFuncValue(f)
}

func VCBIConstructor(val *reflect.Value, i interface{}) interfaces.VariadicValue {
	f, _ := val.Interface().(VCBI) // If we call this func, we know we have the right type
	return VariadicFuncValue(func(x []string) error { return f(x, i) })
}

var CallbackWrappers = map[reflect.Type]CallbackWrapper{
	reflect.TypeOf(BCB(nil)):   BCBConstructor,
	reflect.TypeOf(BCBI(nil)):  BCBIConstructor,
	reflect.TypeOf(BCBE(nil)):  BCBEConstructor,
	reflect.TypeOf(BCBIE(nil)): BCBIEConstructor,
	reflect.TypeOf(CB(nil)):    CBConstructor,
	reflect.TypeOf(CBI(nil)):   CBIConstructor,
}

var VariadicCallbackWrappers = map[reflect.Type]VariadicCallbackWrapper{
	reflect.TypeOf(VCB(nil)):  VCBConstructor,
	reflect.TypeOf(VCBI(nil)): VCBIConstructor,
}

func AsCallback(val *reflect.Value, i interface{}) interfaces.FlagValue {
	if wrap, ok := CallbackWrappers[val.Type()]; ok {
		return wrap(val, i)
	}

	return nil // couldn't convert
}

func AsVariadicCallback(val *reflect.Value, i interface{}) interfaces.VariadicValue {
	if wrap, ok := VariadicCallbackWrappers[val.Type()]; ok {
		return wrap(val, i)
	}

	return nil // couldn't convert
}

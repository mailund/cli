package vals

import (
	"reflect"

	"github.com/mailund/cli/interfaces"
)

// Function callbacks are special values; we don't need to generate them
// since it is easier just to write the little code there is in the first place
type (
	FuncNoValue       func() error         // FuncNoValue wraps functions that do not take parameter arguments
	FuncValue         func(string) error   // FuncValue wraps functions with a single string argument
	VariadicFuncValue func([]string) error // VariadicFuncValue wraps functions with string slice argument
)

func (f FuncNoValue) Set(_ string) error { return f() }  // Set implements the PosValue/FlagValue interface
func (f FuncNoValue) String() string     { return "" }   // String implements the FlagValue interface
func (f FuncNoValue) NoValueFlag() bool  { return true } // NoValueFlag implements the no values interface

// Validate implements the validator interface and checks that callbacks are not nil.
func (f FuncNoValue) Validate(bool) error { // FIXME: the validators are the same, but need generics to join them
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

func (f FuncValue) Set(x string) error { return f(x) }
func (f FuncValue) String() string     { return "" }
func (f FuncValue) Validate(bool) error {
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

func (f VariadicFuncValue) Set(x []string) error { return f(x) }
func (f VariadicFuncValue) Validate(bool) error {
	if f == nil {
		return interfaces.SpecErrorf("callbacks cannot be nil")
	}

	return nil
}

type (
	NCB   = func()
	NCBI  = func(interface{})
	NCBE  = func() error
	NCBIE = func(interface{}) error

	CB  = func(string) error
	CBI = func(string, interface{}) error

	VCB  = func([]string) error
	VCBI = func([]string, interface{}) error

	CallbackWrapper         func(*reflect.Value, interface{}) interfaces.FlagValue
	VariadicCallbackWrapper func(*reflect.Value, interface{}) interfaces.VariadicValue
)

func NCBConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(NCB) // If we call this func, we know we have the right type
	return FuncNoValue(func() error { f(); return nil })
}

func NCBIConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(NCBI) // If we call this func, we know we have the right type
	return FuncNoValue(func() error { f(i); return nil })
}

func NCBEConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(NCBE) // If we call this func, we know we have the right type
	return FuncNoValue(f)
}

func NCBIEConstructor(val *reflect.Value, i interface{}) interfaces.FlagValue {
	f, _ := val.Interface().(NCBIE) // If we call this func, we know we have the right type
	return FuncNoValue(func() error { return f(i) })
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
	reflect.TypeOf(NCB(nil)):   NCBConstructor,
	reflect.TypeOf(NCBI(nil)):  NCBIConstructor,
	reflect.TypeOf(NCBE(nil)):  NCBEConstructor,
	reflect.TypeOf(NCBIE(nil)): NCBIEConstructor,
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

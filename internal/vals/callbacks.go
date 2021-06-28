package vals

import "reflect"

// Function callbacks are special values; we don't need to generate them
// since it is easier just to write the little code there is in the first place

type FuncValue func(string) error

func (f FuncValue) Set(x string) error { return f(x) }
func (f FuncValue) String() string     { return "" }

type VariadicFuncValue func([]string) error

func (f VariadicFuncValue) Set(x []string) error { return f(x) }
func (f VariadicFuncValue) String() string       { return "" }

type (
	CB  = func(string) error
	CBI = func(string, interface{}) error

	VCB  = func([]string) error
	VCBI = func([]string, interface{}) error

	CallbackWrapper         func(*reflect.Value, interface{}) FuncValue
	VariadicCallbackWrapper func(*reflect.Value, interface{}) VariadicFuncValue
)

func CBConstructor(val *reflect.Value, i interface{}) FuncValue {
	return FuncValue(val.Interface().(CB))
}

func CBIConstructor(val *reflect.Value, i interface{}) FuncValue {
	return FuncValue(func(x string) error {
		return val.Interface().(CBI)(x, i)
	})
}

func VCBConstructor(val *reflect.Value, i interface{}) VariadicFuncValue {
	return VariadicFuncValue(val.Interface().(VCB))
}

func VCBIConstructor(val *reflect.Value, i interface{}) VariadicFuncValue {
	return VariadicFuncValue(func(x []string) error {
		return val.Interface().(VCBI)(x, i)
	})
}

var CallbackWrappers = map[reflect.Type]CallbackWrapper{
	reflect.TypeOf(CB(nil)):  CBConstructor,
	reflect.TypeOf(CBI(nil)): CBIConstructor,
}

var VariadicCallbackWrappers = map[reflect.Type]VariadicCallbackWrapper{
	reflect.TypeOf(VCB(nil)):  VCBConstructor,
	reflect.TypeOf(VCBI(nil)): VCBIConstructor,
}

func AsCallback(val *reflect.Value, i interface{}) FuncValue {
	if wrap, ok := CallbackWrappers[val.Type()]; ok {
		return wrap(val, i)
	}

	return nil // couldn't convert
}

func AsVariadicCallback(val *reflect.Value, i interface{}) VariadicFuncValue {
	if wrap, ok := VariadicCallbackWrappers[val.Type()]; ok {
		return wrap(val, i)
	}

	return nil // couldn't convert
}

// Code generated by gen/genvals.go DO NOT EDIT.
package vals

import (
	"reflect"
	"strconv"

	"github.com/mailund/cli/inter"
)

type StringValue string

func (val *StringValue) Set(x string) error {
	*val = StringValue(x)
	return nil
}

func (val *StringValue) String() string {
	return string(*val)
}

func StringValueConstructor(val reflect.Value) inter.FlagValue {
	return (*StringValue)(val.Interface().(*string))
}

type VariadicStringValue []string

func (vals *VariadicStringValue) Set(xs []string) error {
	*vals = make([]string, len(xs))

	for i, x := range xs {
		(*vals)[i] = string(x)
	}

	return nil
}

func VariadicStringValueConstructor(val reflect.Value) inter.VariadicValue {
	return (*VariadicStringValue)(val.Interface().(*[]string))
}

type BoolValue bool

func (val *BoolValue) Set(x string) error {
	v, err := strconv.ParseBool(x)
	if err != nil {
		err = inter.ParseErrorf("argument \"%s\" cannot be parsed as bool", x)
	} else {
		*val = BoolValue(v)
	}

	return err
}

func (val *BoolValue) String() string {
	return strconv.FormatBool(bool(*val))
}

func BoolValueConstructor(val reflect.Value) inter.FlagValue {
	return (*BoolValue)(val.Interface().(*bool))
}

type VariadicBoolValue []bool

func (vals *VariadicBoolValue) Set(xs []string) error {
	*vals = make([]bool, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseBool(x)
		if err != nil {
			return inter.ParseErrorf("cannot parse '%s' as bool", x)
		}

		(*vals)[i] = bool(val)
	}

	return nil
}

func VariadicBoolValueConstructor(val reflect.Value) inter.VariadicValue {
	return (*VariadicBoolValue)(val.Interface().(*[]bool))
}

type IntValue int

func (val *IntValue) Set(x string) error {
	v, err := strconv.ParseInt(x, 0, strconv.IntSize)
	if err != nil {
		err = inter.ParseErrorf("argument \"%s\" cannot be parsed as int", x)
	} else {
		*val = IntValue(v)
	}

	return err
}

func (val *IntValue) String() string {
	return strconv.Itoa(int(*val))
}

func IntValueConstructor(val reflect.Value) inter.FlagValue {
	return (*IntValue)(val.Interface().(*int))
}

type VariadicIntValue []int

func (vals *VariadicIntValue) Set(xs []string) error {
	*vals = make([]int, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseInt(x, 0, strconv.IntSize)
		if err != nil {
			return inter.ParseErrorf("cannot parse '%s' as int", x)
		}

		(*vals)[i] = int(val)
	}

	return nil
}

func VariadicIntValueConstructor(val reflect.Value) inter.VariadicValue {
	return (*VariadicIntValue)(val.Interface().(*[]int))
}

type Float64Value float64

func (val *Float64Value) Set(x string) error {
	v, err := strconv.ParseFloat(x, 64)
	if err != nil {
		err = inter.ParseErrorf("argument \"%s\" cannot be parsed as float64", x)
	} else {
		*val = Float64Value(v)
	}

	return err
}

func (val *Float64Value) String() string {
	return strconv.FormatFloat(float64(*val), 'g', -1, 64)
}

func Float64ValueConstructor(val reflect.Value) inter.FlagValue {
	return (*Float64Value)(val.Interface().(*float64))
}

type VariadicFloat64Value []float64

func (vals *VariadicFloat64Value) Set(xs []string) error {
	*vals = make([]float64, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return inter.ParseErrorf("cannot parse '%s' as float64", x)
		}

		(*vals)[i] = float64(val)
	}

	return nil
}

func VariadicFloat64ValueConstructor(val reflect.Value) inter.VariadicValue {
	return (*VariadicFloat64Value)(val.Interface().(*[]float64))
}

func init() {
	ValsConstructors[reflect.TypeOf((*string)(nil))] = StringValueConstructor
	VarValsConstructors[reflect.TypeOf((*[]string)(nil))] = VariadicStringValueConstructor
	ValsConstructors[reflect.TypeOf((*bool)(nil))] = BoolValueConstructor
	VarValsConstructors[reflect.TypeOf((*[]bool)(nil))] = VariadicBoolValueConstructor
	ValsConstructors[reflect.TypeOf((*int)(nil))] = IntValueConstructor
	VarValsConstructors[reflect.TypeOf((*[]int)(nil))] = VariadicIntValueConstructor
	ValsConstructors[reflect.TypeOf((*float64)(nil))] = Float64ValueConstructor
	VarValsConstructors[reflect.TypeOf((*[]float64)(nil))] = VariadicFloat64ValueConstructor
}

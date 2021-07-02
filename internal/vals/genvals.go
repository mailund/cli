// Code generated by gen/genvals.go DO NOT EDIT.
package vals

import (
	"reflect"
	"strconv"

	"github.com/mailund/cli/interfaces"
)

type StringValue string

func (val *StringValue) Set(x string) error {
	*val = StringValue(x)
	return nil
}

func (val *StringValue) String() string {
	return string(*val)
}

func (val *StringValue) FlagValueDescription() string {
	return "string"
}

func StringValueConstructor(val reflect.Value) interfaces.FlagValue {
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

func (val *VariadicStringValue) FlagValueDescription() string {
	return "string(s)"
}


func VariadicStringValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicStringValue)(val.Interface().(*[]string))
}

type BoolValue bool

func (val *BoolValue) Set(x string) error {
	v, err := strconv.ParseBool(x)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as bool", x)
	} else {
		*val = BoolValue(v)
	}

	return err
}

func (val *BoolValue) String() string {
	return strconv.FormatBool(bool(*val))
}

func (val *BoolValue) FlagValueDescription() string {
	return "boolean"
}

func BoolValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*BoolValue)(val.Interface().(*bool))
}

type VariadicBoolValue []bool

func (vals *VariadicBoolValue) Set(xs []string) error {
	*vals = make([]bool, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseBool(x)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as bool", x)
		}

		(*vals)[i] = bool(val)
	}

	return nil
}

func (val *VariadicBoolValue) FlagValueDescription() string {
	return "boolean(s)"
}


func VariadicBoolValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicBoolValue)(val.Interface().(*[]bool))
}

type IntValue int

func (val *IntValue) Set(x string) error {
	v, err := strconv.ParseInt(x, 0, strconv.IntSize)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as int", x)
	} else {
		*val = IntValue(v)
	}

	return err
}

func (val *IntValue) String() string {
	return strconv.Itoa(int(*val))
}

func (val *IntValue) FlagValueDescription() string {
	return "integer"
}

func IntValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*IntValue)(val.Interface().(*int))
}

type VariadicIntValue []int

func (vals *VariadicIntValue) Set(xs []string) error {
	*vals = make([]int, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseInt(x, 0, strconv.IntSize)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as int", x)
		}

		(*vals)[i] = int(val)
	}

	return nil
}

func (val *VariadicIntValue) FlagValueDescription() string {
	return "integer(s)"
}


func VariadicIntValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicIntValue)(val.Interface().(*[]int))
}

type Int8Value int8

func (val *Int8Value) Set(x string) error {
	v, err := strconv.ParseInt(x, 0, 8)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as int8", x)
	} else {
		*val = Int8Value(v)
	}

	return err
}

func (val *Int8Value) String() string {
	return strconv.Itoa(int(*val))
}

func (val *Int8Value) FlagValueDescription() string {
	return "integer"
}

func Int8ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Int8Value)(val.Interface().(*int8))
}

type VariadicInt8Value []int8

func (vals *VariadicInt8Value) Set(xs []string) error {
	*vals = make([]int8, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseInt(x, 0, 8)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as int8", x)
		}

		(*vals)[i] = int8(val)
	}

	return nil
}

func (val *VariadicInt8Value) FlagValueDescription() string {
	return "integer(s)"
}


func VariadicInt8ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicInt8Value)(val.Interface().(*[]int8))
}

type Int16Value int16

func (val *Int16Value) Set(x string) error {
	v, err := strconv.ParseInt(x, 0, 16)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as int16", x)
	} else {
		*val = Int16Value(v)
	}

	return err
}

func (val *Int16Value) String() string {
	return strconv.Itoa(int(*val))
}

func (val *Int16Value) FlagValueDescription() string {
	return "integer"
}

func Int16ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Int16Value)(val.Interface().(*int16))
}

type VariadicInt16Value []int16

func (vals *VariadicInt16Value) Set(xs []string) error {
	*vals = make([]int16, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseInt(x, 0, 16)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as int16", x)
		}

		(*vals)[i] = int16(val)
	}

	return nil
}

func (val *VariadicInt16Value) FlagValueDescription() string {
	return "integer(s)"
}


func VariadicInt16ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicInt16Value)(val.Interface().(*[]int16))
}

type Int32Value int32

func (val *Int32Value) Set(x string) error {
	v, err := strconv.ParseInt(x, 0, 32)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as int32", x)
	} else {
		*val = Int32Value(v)
	}

	return err
}

func (val *Int32Value) String() string {
	return strconv.Itoa(int(*val))
}

func (val *Int32Value) FlagValueDescription() string {
	return "integer"
}

func Int32ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Int32Value)(val.Interface().(*int32))
}

type VariadicInt32Value []int32

func (vals *VariadicInt32Value) Set(xs []string) error {
	*vals = make([]int32, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseInt(x, 0, 32)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as int32", x)
		}

		(*vals)[i] = int32(val)
	}

	return nil
}

func (val *VariadicInt32Value) FlagValueDescription() string {
	return "integer(s)"
}


func VariadicInt32ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicInt32Value)(val.Interface().(*[]int32))
}

type Int64Value int64

func (val *Int64Value) Set(x string) error {
	v, err := strconv.ParseInt(x, 0, 64)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as int64", x)
	} else {
		*val = Int64Value(v)
	}

	return err
}

func (val *Int64Value) String() string {
	return strconv.Itoa(int(*val))
}

func (val *Int64Value) FlagValueDescription() string {
	return "integer"
}

func Int64ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Int64Value)(val.Interface().(*int64))
}

type VariadicInt64Value []int64

func (vals *VariadicInt64Value) Set(xs []string) error {
	*vals = make([]int64, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseInt(x, 0, 64)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as int64", x)
		}

		(*vals)[i] = int64(val)
	}

	return nil
}

func (val *VariadicInt64Value) FlagValueDescription() string {
	return "integer(s)"
}


func VariadicInt64ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicInt64Value)(val.Interface().(*[]int64))
}

type UintValue uint

func (val *UintValue) Set(x string) error {
	v, err := strconv.ParseUint(x, 0, strconv.IntSize)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as uint", x)
	} else {
		*val = UintValue(v)
	}

	return err
}

func (val *UintValue) String() string {
	return strconv.FormatUint(uint64(*val), 10)
}

func (val *UintValue) FlagValueDescription() string {
	return "unsigned integer"
}

func UintValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*UintValue)(val.Interface().(*uint))
}

type VariadicUintValue []uint

func (vals *VariadicUintValue) Set(xs []string) error {
	*vals = make([]uint, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseUint(x, 0, strconv.IntSize)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as uint", x)
		}

		(*vals)[i] = uint(val)
	}

	return nil
}

func (val *VariadicUintValue) FlagValueDescription() string {
	return "unsigned integer(s)"
}


func VariadicUintValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicUintValue)(val.Interface().(*[]uint))
}

type Uint8Value uint8

func (val *Uint8Value) Set(x string) error {
	v, err := strconv.ParseUint(x, 0, 8)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as uint8", x)
	} else {
		*val = Uint8Value(v)
	}

	return err
}

func (val *Uint8Value) String() string {
	return strconv.FormatUint(uint64(*val), 10)
}

func (val *Uint8Value) FlagValueDescription() string {
	return "unsigned integer"
}

func Uint8ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Uint8Value)(val.Interface().(*uint8))
}

type VariadicUint8Value []uint8

func (vals *VariadicUint8Value) Set(xs []string) error {
	*vals = make([]uint8, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseUint(x, 0, 8)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as uint8", x)
		}

		(*vals)[i] = uint8(val)
	}

	return nil
}

func (val *VariadicUint8Value) FlagValueDescription() string {
	return "unsigned integer(s)"
}


func VariadicUint8ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicUint8Value)(val.Interface().(*[]uint8))
}

type Uint16Value uint16

func (val *Uint16Value) Set(x string) error {
	v, err := strconv.ParseUint(x, 0, 16)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as uint16", x)
	} else {
		*val = Uint16Value(v)
	}

	return err
}

func (val *Uint16Value) String() string {
	return strconv.FormatUint(uint64(*val), 10)
}

func (val *Uint16Value) FlagValueDescription() string {
	return "unsigned integer"
}

func Uint16ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Uint16Value)(val.Interface().(*uint16))
}

type VariadicUint16Value []uint16

func (vals *VariadicUint16Value) Set(xs []string) error {
	*vals = make([]uint16, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseUint(x, 0, 16)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as uint16", x)
		}

		(*vals)[i] = uint16(val)
	}

	return nil
}

func (val *VariadicUint16Value) FlagValueDescription() string {
	return "unsigned integer(s)"
}


func VariadicUint16ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicUint16Value)(val.Interface().(*[]uint16))
}

type Uint32Value uint32

func (val *Uint32Value) Set(x string) error {
	v, err := strconv.ParseUint(x, 0, 32)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as uint32", x)
	} else {
		*val = Uint32Value(v)
	}

	return err
}

func (val *Uint32Value) String() string {
	return strconv.FormatUint(uint64(*val), 10)
}

func (val *Uint32Value) FlagValueDescription() string {
	return "unsigned integer"
}

func Uint32ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Uint32Value)(val.Interface().(*uint32))
}

type VariadicUint32Value []uint32

func (vals *VariadicUint32Value) Set(xs []string) error {
	*vals = make([]uint32, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseUint(x, 0, 32)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as uint32", x)
		}

		(*vals)[i] = uint32(val)
	}

	return nil
}

func (val *VariadicUint32Value) FlagValueDescription() string {
	return "unsigned integer(s)"
}


func VariadicUint32ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicUint32Value)(val.Interface().(*[]uint32))
}

type Uint64Value uint64

func (val *Uint64Value) Set(x string) error {
	v, err := strconv.ParseUint(x, 0, 64)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as uint64", x)
	} else {
		*val = Uint64Value(v)
	}

	return err
}

func (val *Uint64Value) String() string {
	return strconv.FormatUint(uint64(*val), 10)
}

func (val *Uint64Value) FlagValueDescription() string {
	return "unsigned integer"
}

func Uint64ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Uint64Value)(val.Interface().(*uint64))
}

type VariadicUint64Value []uint64

func (vals *VariadicUint64Value) Set(xs []string) error {
	*vals = make([]uint64, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseUint(x, 0, 64)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as uint64", x)
		}

		(*vals)[i] = uint64(val)
	}

	return nil
}

func (val *VariadicUint64Value) FlagValueDescription() string {
	return "unsigned integer(s)"
}


func VariadicUint64ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicUint64Value)(val.Interface().(*[]uint64))
}

type Float32Value float32

func (val *Float32Value) Set(x string) error {
	v, err := strconv.ParseFloat(x, 32)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as float32", x)
	} else {
		*val = Float32Value(v)
	}

	return err
}

func (val *Float32Value) String() string {
	return strconv.FormatFloat(float64(*val), 'g', -1, 32)
}

func (val *Float32Value) FlagValueDescription() string {
	return "floating point number"
}

func Float32ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Float32Value)(val.Interface().(*float32))
}

type VariadicFloat32Value []float32

func (vals *VariadicFloat32Value) Set(xs []string) error {
	*vals = make([]float32, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseFloat(x, 32)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as float32", x)
		}

		(*vals)[i] = float32(val)
	}

	return nil
}

func (val *VariadicFloat32Value) FlagValueDescription() string {
	return "floating point number(s)"
}


func VariadicFloat32ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicFloat32Value)(val.Interface().(*[]float32))
}

type Float64Value float64

func (val *Float64Value) Set(x string) error {
	v, err := strconv.ParseFloat(x, 64)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as float64", x)
	} else {
		*val = Float64Value(v)
	}

	return err
}

func (val *Float64Value) String() string {
	return strconv.FormatFloat(float64(*val), 'g', -1, 64)
}

func (val *Float64Value) FlagValueDescription() string {
	return "floating point number"
}

func Float64ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Float64Value)(val.Interface().(*float64))
}

type VariadicFloat64Value []float64

func (vals *VariadicFloat64Value) Set(xs []string) error {
	*vals = make([]float64, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as float64", x)
		}

		(*vals)[i] = float64(val)
	}

	return nil
}

func (val *VariadicFloat64Value) FlagValueDescription() string {
	return "floating point number(s)"
}


func VariadicFloat64ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicFloat64Value)(val.Interface().(*[]float64))
}

type Complex64Value complex64

func (val *Complex64Value) Set(x string) error {
	v, err := strconv.ParseComplex(x, 64)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as complex64", x)
	} else {
		*val = Complex64Value(v)
	}

	return err
}

func (val *Complex64Value) String() string {
	return strconv.FormatComplex(complex128(*val), 'g', -1, 64)
}

func (val *Complex64Value) FlagValueDescription() string {
	return "complex number"
}

func Complex64ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Complex64Value)(val.Interface().(*complex64))
}

type VariadicComplex64Value []complex64

func (vals *VariadicComplex64Value) Set(xs []string) error {
	*vals = make([]complex64, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseComplex(x, 64)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as complex64", x)
		}

		(*vals)[i] = complex64(val)
	}

	return nil
}

func (val *VariadicComplex64Value) FlagValueDescription() string {
	return "complex number(s)"
}


func VariadicComplex64ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicComplex64Value)(val.Interface().(*[]complex64))
}

type Complex128Value complex128

func (val *Complex128Value) Set(x string) error {
	v, err := strconv.ParseComplex(x, 128)
	if err != nil {
		err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as complex128", x)
	} else {
		*val = Complex128Value(v)
	}

	return err
}

func (val *Complex128Value) String() string {
	return strconv.FormatComplex(complex128(*val), 'g', -1, 128)
}

func (val *Complex128Value) FlagValueDescription() string {
	return "complex number"
}

func Complex128ValueConstructor(val reflect.Value) interfaces.FlagValue {
	return (*Complex128Value)(val.Interface().(*complex128))
}

type VariadicComplex128Value []complex128

func (vals *VariadicComplex128Value) Set(xs []string) error {
	*vals = make([]complex128, len(xs))

	for i, x := range xs {
		val, err := strconv.ParseComplex(x, 128)
		if err != nil {
			return interfaces.ParseErrorf("cannot parse '%s' as complex128", x)
		}

		(*vals)[i] = complex128(val)
	}

	return nil
}

func (val *VariadicComplex128Value) FlagValueDescription() string {
	return "complex number(s)"
}


func VariadicComplex128ValueConstructor(val reflect.Value) interfaces.VariadicValue {
	return (*VariadicComplex128Value)(val.Interface().(*[]complex128))
}

func init() {
	valsConstructors[reflect.TypeOf((*string)(nil))] = StringValueConstructor
	varValsConstructors[reflect.TypeOf((*[]string)(nil))] = VariadicStringValueConstructor
	valsConstructors[reflect.TypeOf((*bool)(nil))] = BoolValueConstructor
	varValsConstructors[reflect.TypeOf((*[]bool)(nil))] = VariadicBoolValueConstructor
	valsConstructors[reflect.TypeOf((*int)(nil))] = IntValueConstructor
	varValsConstructors[reflect.TypeOf((*[]int)(nil))] = VariadicIntValueConstructor
	valsConstructors[reflect.TypeOf((*int8)(nil))] = Int8ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]int8)(nil))] = VariadicInt8ValueConstructor
	valsConstructors[reflect.TypeOf((*int16)(nil))] = Int16ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]int16)(nil))] = VariadicInt16ValueConstructor
	valsConstructors[reflect.TypeOf((*int32)(nil))] = Int32ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]int32)(nil))] = VariadicInt32ValueConstructor
	valsConstructors[reflect.TypeOf((*int64)(nil))] = Int64ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]int64)(nil))] = VariadicInt64ValueConstructor
	valsConstructors[reflect.TypeOf((*uint)(nil))] = UintValueConstructor
	varValsConstructors[reflect.TypeOf((*[]uint)(nil))] = VariadicUintValueConstructor
	valsConstructors[reflect.TypeOf((*uint8)(nil))] = Uint8ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]uint8)(nil))] = VariadicUint8ValueConstructor
	valsConstructors[reflect.TypeOf((*uint16)(nil))] = Uint16ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]uint16)(nil))] = VariadicUint16ValueConstructor
	valsConstructors[reflect.TypeOf((*uint32)(nil))] = Uint32ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]uint32)(nil))] = VariadicUint32ValueConstructor
	valsConstructors[reflect.TypeOf((*uint64)(nil))] = Uint64ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]uint64)(nil))] = VariadicUint64ValueConstructor
	valsConstructors[reflect.TypeOf((*float32)(nil))] = Float32ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]float32)(nil))] = VariadicFloat32ValueConstructor
	valsConstructors[reflect.TypeOf((*float64)(nil))] = Float64ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]float64)(nil))] = VariadicFloat64ValueConstructor
	valsConstructors[reflect.TypeOf((*complex64)(nil))] = Complex64ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]complex64)(nil))] = VariadicComplex64ValueConstructor
	valsConstructors[reflect.TypeOf((*complex128)(nil))] = Complex128ValueConstructor
	varValsConstructors[reflect.TypeOf((*[]complex128)(nil))] = VariadicComplex128ValueConstructor
}

package cli

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mailund/cli/internal/params"
)

// SpecError is the error type returned if there are problems with a
// specification
type SpecError struct {
	Message string
}

// SpecErrorf creates a SpecError from a format string and arguments
func SpecErrorf(format string, args ...interface{}) *SpecError {
	return &SpecError{fmt.Sprintf(format, args...)}
}

// Error returns a string representation of a SpecError, implementing
// the error interface.
func (err *SpecError) Error() string {
	return err.Message
}

// FIXME: there must be a better way of getting the reflection type
// of a function type...
var uglyHack func(string) error // just a way to get the signature
var callbackSignature = reflect.TypeOf(uglyHack)

func setFlag(f *flag.FlagSet, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	switch tfield.Type.Kind() {
	case reflect.Bool:
		f.BoolVar(vfield.Addr().Interface().(*bool), name, vfield.Bool(), tfield.Tag.Get("descr"))

	case reflect.Int:
		f.IntVar(vfield.Addr().Interface().(*int), name, int(vfield.Int()), tfield.Tag.Get("descr"))

	case reflect.Float64:
		f.Float64Var(vfield.Addr().Interface().(*float64), name, vfield.Float(), tfield.Tag.Get("descr"))

	case reflect.String:
		f.StringVar(vfield.Addr().Interface().(*string), name, vfield.String(), tfield.Tag.Get("descr"))

	case reflect.Func:
		if tfield.Type != callbackSignature {
			return SpecErrorf("callbacks must have signature func(string) error")
		}

		if vfield.IsNil() {
			return SpecErrorf("callbacks cannot be nil")
		}

		f.Func(name, tfield.Tag.Get("descr"), vfield.Interface().(func(string) error))

	default:
		return SpecErrorf("unsupported type for flag %s: %q", name, tfield.Type.Kind())
	}

	return nil
}

func setParam(p *params.ParamSet, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	switch tfield.Type.Kind() {
	case reflect.Bool:
		p.BoolVar(vfield.Addr().Interface().(*bool), name, tfield.Tag.Get("descr"))

	case reflect.Int:
		p.IntVar(vfield.Addr().Interface().(*int), name, tfield.Tag.Get("descr"))

	case reflect.Float64:
		p.FloatVar(vfield.Addr().Interface().(*float64), name, tfield.Tag.Get("descr"))

	case reflect.String:
		p.StringVar(vfield.Addr().Interface().(*string), name, tfield.Tag.Get("descr"))

	case reflect.Func:
		if tfield.Type != callbackSignature {
			return SpecErrorf("callbacks must have signature func(string) error")
		}

		if vfield.IsNil() {
			return SpecErrorf("callbacks cannot be nil")
		}

		p.Func(name, tfield.Tag.Get("descr"), vfield.Interface().(func(string) error))

	default:
		return SpecErrorf("unsupported type for parameter %s: %q", name, tfield.Type.Kind())
	}

	return nil
}

func setVariadicParam(p *params.ParamSet, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	var (
		min int
		err error
	)

	if minTag := tfield.Tag.Get("min"); minTag == "" {
		min = 0
	} else if min, err = strconv.Atoi(tfield.Tag.Get("min")); err != nil {
		return SpecErrorf("unexpected min value for variadic parameter %s: %s", name, minTag)
	}

	switch tfield.Type.Elem().Kind() {
	case reflect.Bool:
		p.VariadicBoolVar(vfield.Addr().Interface().(*[]bool), name, tfield.Tag.Get("descr"), min)

	case reflect.Int:
		p.VariadicIntVar(vfield.Addr().Interface().(*[]int), name, tfield.Tag.Get("descr"), min)

	case reflect.Float64:
		p.VariadicFloatVar(vfield.Addr().Interface().(*[]float64), name, tfield.Tag.Get("descr"), min)

	case reflect.String:
		p.VariadicStringVar(vfield.Addr().Interface().(*[]string), name, tfield.Tag.Get("descr"), min)

	default:
		return SpecErrorf("unsupported slice type for parameter %s: %q", name, tfield.Type.Elem().Kind())
	}

	return nil
}

func connectSpecsFlagsAndParams(f *flag.FlagSet, p *params.ParamSet, argv interface{}, allowVariadic bool) error {
	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()
	seenVariadic := false

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		if name, isFlag := tfield.Tag.Lookup("flag"); isFlag {
			if err := setFlag(f, name, &tfield, &vfield); err != nil {
				return err
			}
		}

		if name, isPos := tfield.Tag.Lookup("pos"); isPos {
			if tfield.Type.Kind() == reflect.Slice {
				if !allowVariadic {
					return SpecErrorf("a command with subcommands cannot have variadic parameters")
				}

				if seenVariadic {
					return SpecErrorf("a command spec cannot contain more than one variadic parameter")
				}

				seenVariadic = true

				if err := setVariadicParam(p, name, &tfield, &vfield); err != nil {
					return err
				}
			} else if err := setParam(p, name, &tfield, &vfield); err != nil {
				return err
			}
		}
	}

	return nil
}

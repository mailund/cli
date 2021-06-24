package cli

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mailund/cli/params"
)

type SpecError struct {
	Message string
}

func SpecErrorf(format string, args ...interface{}) *SpecError {
	return &SpecError{fmt.Sprintf(format, args...)}
}

func (err *SpecError) Error() string {
	return err.Message
}

func prepareSpecsFlags(f *flag.FlagSet, argv interface{}) error {
	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		if name, isFlag := tfield.Tag.Lookup("flag"); isFlag {
			switch tfield.Type.Kind() {
			case reflect.Bool:
				f.BoolVar(vfield.Addr().Interface().(*bool), name, vfield.Bool(), tfield.Tag.Get("descr"))

			case reflect.Int:
				f.IntVar(vfield.Addr().Interface().(*int), name, int(vfield.Int()), tfield.Tag.Get("descr"))

			case reflect.Float64:
				f.Float64Var(vfield.Addr().Interface().(*float64), name, vfield.Float(), tfield.Tag.Get("descr"))

			case reflect.String:
				f.StringVar(vfield.Addr().Interface().(*string), name, vfield.String(), tfield.Tag.Get("descr"))

			default:
				return SpecErrorf("unsupported type for flag %s: %q", name, tfield.Type.Kind())
			}
		}
	}

	return nil
}

// TODO: check number of variadic
// TODO: add Func arguments.
// TODO: connect specs with commands; no variadic for commands with subcommands

func prepareSpecsParams(p *params.ParamSet, argv interface{}) error {
	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		if name, isArg := tfield.Tag.Lookup("arg"); isArg {
			switch tfield.Type.Kind() {
			case reflect.Bool:
				p.BoolVar(vfield.Addr().Interface().(*bool), name, tfield.Tag.Get("descr"))

			case reflect.Int:
				p.IntVar(vfield.Addr().Interface().(*int), name, tfield.Tag.Get("descr"))

			case reflect.Float64:
				p.FloatVar(vfield.Addr().Interface().(*float64), name, tfield.Tag.Get("descr"))

			case reflect.String:
				p.StringVar(vfield.Addr().Interface().(*string), name, tfield.Tag.Get("descr"))

			case reflect.Slice:
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

			default:
				return SpecErrorf("unsupported type for parameter %s: %q", name, tfield.Type.Kind())
			}
		}
	}

	return nil
}

func prepareSpecs(f *flag.FlagSet, p *params.ParamSet, argv interface{}) error {
	if err := prepareSpecsFlags(f, argv); err != nil {
		return err
	}

	return prepareSpecsParams(p, argv)
}

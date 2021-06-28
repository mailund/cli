package cli

import (
	"reflect"
	"strconv"

	"github.com/mailund/cli/inter"
	"github.com/mailund/cli/internal/vals"
)

func setFlag(cmd *Command, argv interface{}, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	// bool flags are handled differently from values, so we insert those explicitly as books...
	if tfield.Type.Kind() == reflect.Bool {
		cmd.flags.BoolVar(vfield.Addr().Interface().(*bool), name, vfield.Bool(), tfield.Tag.Get("descr"))
		return nil
	}

	// Otherwise, rely on our value wrappers...
	if val := vals.AsValue(vfield.Addr()); val != nil {
		cmd.flags.Var(val, name, tfield.Tag.Get("descr"))
		return nil
	}

	if tfield.Type.Kind() == reflect.Func {
		if vfield.IsNil() {
			return inter.SpecErrorf("callbacks cannot be nil")
		}

		if val := vals.AsCallback(vfield, argv); val != nil {
			cmd.flags.Var(val, name, tfield.Tag.Get("descr"))
			return nil
		}

		return inter.SpecErrorf("incorrect signature for callbacks: %q", tfield.Type)
	}

	return inter.SpecErrorf("unsupported type for flag %s: %q", name, tfield.Type.Kind())
}

func setParam(cmd *Command, argv interface{}, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	if val := vals.AsValue(vfield.Addr()); val != nil {
		cmd.params.Var(val, name, tfield.Tag.Get("descr"))
		return nil
	}
	// FIXME: handle variadic

	if tfield.Type.Kind() == reflect.Func {
		if vfield.IsNil() {
			return inter.SpecErrorf("callbacks cannot be nil")
		}

		if val := vals.AsCallback(vfield, argv); val != nil {
			cmd.params.Var(val, name, tfield.Tag.Get("descr"))
			return nil
		}
		// FIXME: handle variadic

		return inter.SpecErrorf("incorrect signature for callbacks: %q", tfield.Type)
	}

	return inter.SpecErrorf("unsupported type for parameter %s: %q", name, tfield.Type.Kind())
}

func setVariadicParam(cmd *Command, _ interface{}, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	var (
		min int
		err error
	)

	if minTag := tfield.Tag.Get("min"); minTag == "" {
		min = 0
	} else if min, err = strconv.Atoi(tfield.Tag.Get("min")); err != nil {
		return inter.SpecErrorf("unexpected min value for variadic parameter %s: %s", name, minTag)
	}

	switch tfield.Type.Elem().Kind() {
	case reflect.Bool:
		cmd.params.VariadicBoolVar(vfield.Addr().Interface().(*[]bool), name, tfield.Tag.Get("descr"), min)

	case reflect.Int:
		cmd.params.VariadicIntVar(vfield.Addr().Interface().(*[]int), name, tfield.Tag.Get("descr"), min)

	case reflect.Float64:
		cmd.params.VariadicFloatVar(vfield.Addr().Interface().(*[]float64), name, tfield.Tag.Get("descr"), min)

	case reflect.String:
		cmd.params.VariadicStringVar(vfield.Addr().Interface().(*[]string), name, tfield.Tag.Get("descr"), min)

	default:
		return inter.SpecErrorf("unsupported slice type for parameter %s: %q", name, tfield.Type.Elem().Kind())
	}

	return nil
}

func connectSpecsFlagsAndParams(cmd *Command, argv interface{}) error {
	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()
	seenVariadic := false

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		if name, isFlag := tfield.Tag.Lookup("flag"); isFlag {
			if err := setFlag(cmd, argv, name, &tfield, &vfield); err != nil {
				return err
			}
		}

		if name, isPos := tfield.Tag.Lookup("pos"); isPos {
			if tfield.Type.Kind() == reflect.Slice {
				if len(cmd.Subcommands) > 0 {
					return inter.SpecErrorf("a command with subcommands cannot have variadic parameters")
				}

				if seenVariadic {
					return inter.SpecErrorf("a command spec cannot contain more than one variadic parameter")
				}

				seenVariadic = true

				if err := setVariadicParam(cmd, argv, name, &tfield, &vfield); err != nil {
					return err
				}
			} else if err := setParam(cmd, argv, name, &tfield, &vfield); err != nil {
				return err
			}
		}
	}

	return nil
}

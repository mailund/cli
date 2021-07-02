package cli

import (
	"reflect"
	"strconv"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/vals"
)

func setFlag(cmd *Command, argv interface{}, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	val := vals.AsFlagValue(vfield.Addr())
	if val == nil {
		val = vals.AsCallback(vfield, argv)
	}

	if val != nil { // We have a value
		short := tfield.Tag.Get("short")
		if short == "" && len(name) == 1 {
			short = name
		}

		return cmd.flags.Var(val, name, short, tfield.Tag.Get("descr"))
	}

	// report appropriate error...
	if tfield.Type.Kind() == reflect.Func {
		return interfaces.SpecErrorf("incorrect signature for callbacks: %q", tfield.Type)
	}

	return interfaces.SpecErrorf("unsupported type for flag %s: %q", name, tfield.Type.Kind())
}

func setVariadic(cmd *Command, name string, val interfaces.VariadicValue, tfield *reflect.StructField) error {
	if len(cmd.Subcommands) > 0 {
		return interfaces.SpecErrorf("a command with subcommands cannot have variadic parameters")
	}

	if cmd.params.Variadic() != nil {
		return interfaces.SpecErrorf("a command spec cannot contain more than one variadic parameter")
	}

	var (
		min int
		err error
	)

	if minTag := tfield.Tag.Get("min"); minTag == "" {
		min = 0
	} else if min, err = strconv.Atoi(tfield.Tag.Get("min")); err != nil {
		return interfaces.SpecErrorf("unexpected min value for variadic parameter %s: %s", name, minTag)
	}

	cmd.params.VariadicVar(val, name, tfield.Tag.Get("descr"), min)

	return nil
}

func setParam(cmd *Command, argv interface{}, name string, tfield *reflect.StructField, vfield *reflect.Value) error {
	// first, try normal value or callback
	val := vals.AsPosValue(vfield.Addr())
	if val == nil {
		val = vals.AsCallback(vfield, argv)
	}

	if val != nil {
		// we have a value...
		cmd.params.Var(val, name, tfield.Tag.Get("descr"))
		return nil
	}

	// then try variadics...
	if val := vals.AsVariadicValue(vfield.Addr()); val != nil {
		return setVariadic(cmd, name, val, tfield)
	}

	if val := vals.AsVariadicCallback(vfield, argv); val != nil {
		return setVariadic(cmd, name, val, tfield)
	}

	// nothing worked, so we report an appropriate error
	if tfield.Type.Kind() == reflect.Func {
		return interfaces.SpecErrorf("incorrect signature for callbacks: %q", tfield.Type)
	}

	return interfaces.SpecErrorf("unsupported type for parameter %s: %q", name, tfield.Type.Kind())
}

func connectSpecsFlagsAndParams(cmd *Command, argv interface{}) error {
	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		if name, isFlag := tfield.Tag.Lookup("flag"); isFlag {
			if err := setFlag(cmd, argv, name, &tfield, &vfield); err != nil {
				return err
			}
		}

		if name, isPos := tfield.Tag.Lookup("pos"); isPos {
			if err := setParam(cmd, argv, name, &tfield, &vfield); err != nil {
				return err
			}
		}
	}

	return validateFlagsAndParams(cmd)
}

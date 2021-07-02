package cli

import (
	"github.com/mailund/cli/interfaces"
)

func prepare(i interface{}) error {
	if v, ok := i.(interfaces.Prepare); ok {
		return v.PrepareValue()
	}

	return nil
}

func prepareFlagsAndParams(cmd *Command) error {
	var err error

	for i := 0; i < cmd.flags.NFlags(); i++ {
		if err = prepare(cmd.flags.Flag(i).Value); err == nil {
			continue
		}

		// we have an error, but need a better error message
		f := cmd.flags.Flag(i)
		flagname := ""

		if f.Short != "" {
			flagname += "-" + f.Short
		}

		if f.Long != "" {
			if f.Short != "" {
				flagname += ","
			}

			flagname += "--" + f.Long
		}

		return interfaces.ParseErrorf("error in flag %s: %s", flagname, err)
	}

	for i := 0; i < cmd.params.NParams(); i++ {
		if err = prepare(cmd.params.Param(i).Value); err != nil {
			return interfaces.ParseErrorf("error in argument %s: %s", cmd.params.Param(i).Name, err)
		}
	}

	if vv := cmd.params.Variadic(); vv != nil {
		if err := prepare(vv.Value); err != nil {
			return interfaces.ParseErrorf("error in argument %s: %s", cmd.params.Variadic().Name, err)
		}
	}

	return nil
}

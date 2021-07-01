package cli

import (
	"github.com/mailund/cli/interfaces"
)

func prepareFlagsAndParams(cmd *Command) error {
	for i := 0; i < cmd.flags.NFlags(); i++ {
		if v, ok := cmd.flags.Flag(i).Value.(interfaces.Prepare); ok {
			if err := v.PrepareValue(); err != nil {
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
		}
	}

	for i := 0; i < cmd.params.NParams(); i++ {
		if v, ok := cmd.params.Param(i).Value.(interfaces.Prepare); ok {
			if err := v.PrepareValue(); err != nil {
				return interfaces.ParseErrorf("error in argument %s: %s", cmd.params.Param(i).Name, err)
			}
		}
	}

	if vv := cmd.params.Variadic(); vv != nil {
		if v, ok := vv.Value.(interfaces.Prepare); ok {
			if err := v.PrepareValue(); err != nil {
				return interfaces.ParseErrorf("error in argument %s: %s", cmd.params.Variadic().Name, err)
			}
		}
	}

	return nil
}

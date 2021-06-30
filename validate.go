package cli

import (
	"github.com/mailund/cli/interfaces"
)

func validateFlagsAndParams(cmd *Command) error {
	for i := 0; i < cmd.flags.NFlags(); i++ {
		if v, ok := cmd.flags.Flag(i).Value.(interfaces.Validator); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	for i := 0; i < cmd.params.NParams(); i++ {
		if v, ok := cmd.params.Param(i).Value.(interfaces.Validator); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	if vv := cmd.params.Variadic(); vv != nil {
		if v, ok := vv.Value.(interfaces.Validator); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

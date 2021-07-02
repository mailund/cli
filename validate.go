package cli

import (
	"github.com/mailund/cli/interfaces"
)

func validate(flag bool, v interface{}) error {
	if val, ok := v.(interfaces.Validator); ok {
		return val.Validate(flag)
	}

	return nil
}

func validateFlagsAndParams(cmd *Command) error {
	for i := 0; i < cmd.flags.NFlags(); i++ {
		if err := validate(true, cmd.flags.Flag(i).Value); err != nil {
			return err
		}
	}

	for i := 0; i < cmd.params.NParams(); i++ {
		if err := validate(false, cmd.params.Param(i).Value); err != nil {
			return err
		}
	}

	if vv := cmd.params.Variadic(); vv != nil {
		return validate(false, vv.Value)
	}

	return nil
}

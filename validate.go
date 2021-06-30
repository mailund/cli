package cli

import "github.com/mailund/cli/interfaces"

func validateFlagsAndParams(cmd *Command) error {
	for i := 0; i < cmd.flags.NFlags(); i++ {
		// FIXME: I can't actually get the fucking flags!
	}

	for i := 0; i < cmd.params.NParams(); i++ {
		if v, ok := cmd.params.Param(i).Value.(interfaces.Validator); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

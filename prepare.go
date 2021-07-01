package cli

import (
	"fmt"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
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

				fmt.Fprintf(cmd.Output(), "error in flag %s: %s", flagname, err)

				if cmd.errf == failure.ExitOnError {
					failure.Failure()
				}

				return err
			}
		}
	}

	for i := 0; i < cmd.params.NParams(); i++ {
		if v, ok := cmd.params.Param(i).Value.(interfaces.Prepare); ok {
			if err := v.PrepareValue(); err != nil {
				fmt.Fprintf(cmd.Output(), "error in argument %s: %s", cmd.params.Param(i).Name, err)

				if cmd.errf == failure.ExitOnError {
					failure.Failure()
				}

				return err
			}
		}
	}

	if vv := cmd.params.Variadic(); vv != nil {
		if v, ok := vv.Value.(interfaces.Prepare); ok {
			if err := v.PrepareValue(); err != nil {
				fmt.Fprintf(cmd.Output(), "error in argument %s: %s", cmd.params.Variadic().Name, err)

				if cmd.errf == failure.ExitOnError {
					failure.Failure()
				}

				return err
			}
		}
	}

	return nil
}
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/mailund/cli/internal/failure"
	"github.com/mailund/cli/params"
)

// Command Wraps a command line (sub)command.
type Command struct {
	// The name of the command, using for subcommands and for documentation
	Name string
	// A short description used for usage for menus of choices of commands.
	ShortDescription string
	// Used for displaying documentation for a command.
	LongDescription string
	// Command line options.
	Flags *flag.FlagSet
	// Command line parameters
	Params *params.ParamSet

	usage func()
	run   func()
	out   io.Writer
}

// Output returns the writer the command will write usage information to.
func (cmd *Command) Output() io.Writer { return cmd.out }

// SetOutput sets the writer that options, arguments, and the command will
// write usage information to.
func (cmd *Command) SetOutput(out io.Writer) {
	cmd.Flags.SetOutput(out)
	cmd.Params.SetOutput(out)
	cmd.out = out
}

// Usage prints usage information about a command.
func (cmd *Command) Usage() { cmd.usage() }

// SetUsage sets the function for printing usage information.
func (cmd *Command) SetUsage(usage func()) {
	cmd.Flags.Usage = usage
	cmd.Params.Usage = usage
	cmd.usage = usage
}

// Run parses options and arguments from args and then executes the
// command.
//
// Parsing errors for either flags or parameters will terminate the program
// with os.Exit(0) for -help options and os.Exit(2) otherwise. If the parsing
// is succesfull, the underlying run callback is executed.
func (cmd *Command) Run(args []string) {
	// Outside of testing, we never get the errors, so we don't
	// propagate them. For testing, it is good to disable the exit
	// in flags and params. The error handling here is only a way of
	// mocking termination and serves no other purpose.
	if err := cmd.Flags.Parse(args); err != nil {
		return
	}

	if err := cmd.Params.Parse(cmd.Flags.Args()); err != nil {
		return
	}

	cmd.run()
}

// NewCommand Create a new command.
//
// Parameters:
//   - name: Name of the command, using when printing usage.
//   - short: Short description of the command, used when printing
//     usage when the command is part of a Menu or when long is empty.
//   - long: Long description of the command, used when printing usage.
//   - init: Used for setting up the underlying callback for executing the
//     command.
//
// The init function will be called with the new Command object, and must
// return a callback that will be executed when the new command is run. The
// purpose of this setup is that init can create variables for the Flags and
// Params fields of the command, and thus put those in its scope. If the
// callback is created in this scope, it has access to them in its closure,
// and when the callback is invoked, the parsed flags and arguments are available
// to it.
func NewCommand(name, short, long string, init func(*Command) func()) *Command {
	cmd := &Command{Name: name,
		ShortDescription: short,
		Flags:            flag.NewFlagSet(name, flag.ExitOnError),
		Params:           params.NewParamSet(name, params.ExitOnError),
		out:              os.Stdout}

	const linewidth = 70
	cmd.LongDescription = wordWrap(long, linewidth)

	// There is always a help command when we parse, but the usage won't
	// show it unless we make it explicit
	cmd.Flags.Bool("help", false, fmt.Sprintf("show help for %s", name))
	cmd.SetUsage(DefaultUsage(cmd))

	// configure according to caller. The caller can modify cmd
	// if necessary, and the init function should create a closure
	// for running the command
	cmd.run = init(cmd)

	return cmd
}

// DefaultUsage creates a default function for printing usage information,
// getting the information to print from the cmd object.
func DefaultUsage(cmd *Command) func() {
	return func() {
		fmt.Fprintf(cmd.Output(),
			"Usage: %s [options] %s\n\n",
			cmd.Name, cmd.Params.ShortUsage())

		if len(cmd.LongDescription) > 0 {
			fmt.Fprintf(cmd.Output(), "%s\n\n", cmd.LongDescription)
		} else {
			fmt.Fprintf(cmd.Output(), "%s\n\n", cmd.ShortDescription)
		}

		// Print options and arguments at the bottom.
		fmt.Fprintf(cmd.Output(), "Options:\n")
		cmd.Flags.PrintDefaults()
		fmt.Fprintf(cmd.Output(), "\n")
		cmd.Params.PrintDefaults()
	}
}

// NewMenu Create a new menu command, i.e. a command that will dispatch to
// sub-commands.
//
// Parameters:
//   - name: Name of the command, used when printing usage.
//   - short: Short description of the command, used when printing
//     usage when the command is part of a Menu or when long is empty.
//   - long: Long description of the command, used when printing usage.
//   - subcmds: The subcommands you can invoke through this menu.
func NewMenu(name, short, long string, subcmds ...*Command) *Command {
	subcommands := map[string]*Command{}
	for _, cmd := range subcmds {
		subcommands[cmd.Name] = cmd
	}

	init := func(cmd *Command) func() {
		var command string

		var cmdArgs []string

		cmd.Params.StringVar(&command, "cmd", "command to run")
		cmd.Params.VariadicStringVar(&cmdArgs, "...", "command arguments", 0)

		usage := func() {
			DefaultUsage(cmd)()
			// add commands to usage...
			fmt.Fprintf(cmd.Output(), "\nCommands:\n")

			subcmdNames := []string{}
			for name := range subcommands {
				subcmdNames = append(subcmdNames, name)
			}

			sort.Strings(subcmdNames)

			for _, name := range subcmdNames {
				fmt.Fprintf(cmd.Output(), "  %s\n\t%s\n", name, subcommands[name].ShortDescription)
			}

			fmt.Fprintf(cmd.Output(), "\n")
		}

		cmd.SetUsage(usage)

		return func() {
			if subcmd, ok := subcommands[command]; ok {
				subcmd.Run(cmdArgs)
			} else {
				fmt.Fprintf(cmd.Output(),
					"'%s' is not a valid command for %s.\n\n", command, name)
				cmd.Usage()
				failure.Failure()
			}
		}
	}

	return NewCommand(name, short, long, init)
}

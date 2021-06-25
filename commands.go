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

// CommandSpec defines a commandline command or subcommand
type CommandSpec struct {
	// Name is the name of the command, using for subcommands and for documentation
	Name string
	// Short is a short description used for usage for menus of choices of commands.
	Short string
	// Long is used for displaying documentation for a command.
	Long string
	// Init is a callback that should create a structure that specifies flags and positional
	// parameters (via reflection) plus default values and any other data the command needs.
	Init func() interface{}
	// Action is called if/when the parser reaches the command. If the commandline has a multi-command
	// path, all actions will be invoked. Commands with subcommands can leave Action as nil to rely on
	// the default behaviour, or handle setup as necessary. Subcommands are handled after the action, if
	// there are any. The argument to Action is the structure returned from Init(), after flags and positional
	// arguments are parsed.
	Action func(interface{})
	// Usage is a callback to print usage information about a command. In most cases, you should leave
	// it undefined and rely on the default usage.
	Usage func()
	// Subcommands holds a list of subcommands.
	Subcommands []*Command
}

// Command wraps a command line (sub)command. It is created from a CommandSpec and is the functional
// part of a command
type Command struct {
	CommandSpec

	flags  *flag.FlagSet
	params *params.ParamSet
	argv   interface{}
	out    io.Writer

	// for subcommands
	subcommands map[string]*Command
	command     string
	cmdArgs     []string
}

// Output returns the writer the command will write usage information to.
func (cmd *Command) Output() io.Writer { return cmd.out }

// SetOutput sets the writer that options, arguments, and the command will
// write usage information to.
func (cmd *Command) SetOutput(out io.Writer) {
	cmd.flags.SetOutput(out)
	cmd.params.SetOutput(out)

	for _, sub := range cmd.Subcommands {
		sub.SetOutput(out)
	}

	cmd.out = out
}

func (cmd *Command) SetErrorFlag(f flag.ErrorHandling) {
	cmd.flags.Init(cmd.Name, f)
	cmd.params.SetFlag(f)

	for _, subcmd := range cmd.Subcommands {
		subcmd.SetErrorFlag(f)
	}
}

// Usage prints usage information about a command.
func (cmd *Command) Usage() { cmd.CommandSpec.Usage() }

// SetUsage sets the function for printing usage information.
// If you only set the Usage field in a spec, it only applies to
// that spec, not the flags or parameters, nor subcommands. You
// should almost always use SetUsage() instead.
func (cmd *Command) SetUsage(usage func()) {
	cmd.flags.Usage = usage
	cmd.params.Usage = usage

	for _, sub := range cmd.Subcommands {
		sub.SetUsage(usage)
	}

	cmd.CommandSpec.Usage = usage
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
	if err := cmd.flags.Parse(args); err != nil {
		return
	}

	if err := cmd.params.Parse(cmd.flags.Args()); err != nil {
		return
	}

	// Invoke the action for this (sub)command
	if cmd.Action != nil {
		cmd.Action(cmd.argv)
	}

	// then, if there are sub-commands, dispatch
	if len(cmd.subcommands) > 0 {
		if subcmd, ok := cmd.subcommands[cmd.command]; ok {
			subcmd.Run(cmd.cmdArgs)
		} else {
			fmt.Fprintf(cmd.Output(),
				"'%s' is not a valid command for %s.\n\n", cmd.command, cmd.Name)
			cmd.Usage()
			failure.Failure()
		}
	}
}

// NewCommand Create a new command. The function returns a new command object or an error.
// Since errors are only possible if the specification is incorrect in some way, you will
// usually want NewCommand, that panics on errors, instead.
//
// The spec defines how the command should behave, though its Init and Action parameters.
// The Init function should return a struct with flags and positional arguments
// annotated so parsed arguments are automatically configured. Then the Action parameter
// will be invoked when the commandline gets to the command.
func NewCommandError(spec CommandSpec) (*Command, error) { //nolint:gocritic // specs are large but only copied when we create a command
	cmd := &Command{
		CommandSpec: spec,
		flags:       flag.NewFlagSet(spec.Name, flag.ExitOnError),
		params:      params.NewParamSet(spec.Name, flag.ExitOnError)}

	if spec.Init != nil {
		cmd.argv = spec.Init()

		if err := connectSpecsFlagsAndParams(cmd.flags, cmd.params, cmd.argv); err != nil {
			return nil, err
		}
	}

	cmd.SetOutput(os.Stdout)

	const linewidth = 70
	cmd.Long = wordWrap(cmd.Long, linewidth)

	// There is always a help command when we parse, but the usage won't
	// show it unless we make it explicit
	cmd.flags.Bool("help", false, fmt.Sprintf("show help for %s", cmd.Name))

	if spec.Usage == nil {
		cmd.SetUsage(DefaultUsage(cmd))
	}

	if len(cmd.Subcommands) > 0 {
		cmd.subcommands = map[string]*Command{}

		for _, sub := range cmd.Subcommands {
			cmd.subcommands[sub.Name] = sub
		}

		cmd.params.StringVar(&cmd.command, "cmd", "sub-command to call")
		cmd.params.VariadicStringVar(&cmd.cmdArgs, "...", "argument for sub-commands", 0)
	}

	return cmd, nil
}

// NewCommand creates a new command from a specification, and panics if there are
// errors in the specification.
//
// The spec defines how the command should behave, though its Init and Action parameters.
// The Init function should return a struct with flags and positional arguments
// annotated so parsed arguments are automatically configured. Then the Action parameter
// will be invoked when the commandline gets to the command.
func NewCommand(spec CommandSpec) *Command { //nolint:gocritic // specs are large but only copied when we create a command
	cmd, err := NewCommandError(spec)
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}

	return cmd
}

// DefaultUsage creates a default function for printing usage information,
// getting the information to print from the cmd object.
func DefaultUsage(cmd *Command) func() {
	return func() {
		fmt.Fprintf(cmd.Output(),
			"Usage: %s [options] %s\n\n",
			cmd.Name, cmd.params.ShortUsage())

		if len(cmd.Long) > 0 {
			fmt.Fprintf(cmd.Output(), "%s\n\n", cmd.Long)
		} else {
			fmt.Fprintf(cmd.Output(), "%s\n\n", cmd.Short)
		}

		// Print options and arguments at the bottom.
		fmt.Fprintf(cmd.Output(), "Options:\n")
		cmd.flags.PrintDefaults()
		fmt.Fprintf(cmd.Output(), "\n")
		cmd.params.PrintDefaults()

		if len(cmd.subcommands) > 0 {
			fmt.Fprintf(cmd.Output(), "\nCommands:\n")

			subcmdNames := []string{}
			for name := range cmd.subcommands {
				subcmdNames = append(subcmdNames, name)
			}

			sort.Strings(subcmdNames)

			for _, name := range subcmdNames {
				fmt.Fprintf(cmd.Output(), "  %s\n\t%s\n", name, cmd.subcommands[name].Short)
			}

			fmt.Fprintf(cmd.Output(), "\n")
		}
	}
}

// NewMenuError Create a new menu command, i.e. a command that will dispatch to
// sub-commands. This is a convinience wrapper around a NewCommand with a spec
// that has subcommands. It returns a new command or an error. It can only error
// if there is an error in the specification, so usually you do not need this function
// but should use NewMenu instead.
//
// Parameters:
//   - name: Name of the command, used when printing usage.
//   - short: Short description of the command, used when printing
//     usage when the command is part of a Menu or when long is empty.
//   - long: Long description of the command, used when printing usage.
//   - subcmds: The subcommands you can invoke through this menu.
func NewMenuError(name, short, long string, subcmds ...*Command) (*Command, error) {
	return NewCommandError(CommandSpec{
		Name:        name,
		Short:       short,
		Long:        long,
		Subcommands: subcmds,
	})
}

// NewMenu creates a menu command and panics if there are errors.
// This is a convinience wrapper around a NewCommand with a spec
// that has subcommands. It returns a new command or an error. It can only error
// if there is an error in the specification, so usually you do not need this function
// but should use NewMenu instead.
//
// Parameters:
//   - name: Name of the command, used when printing usage.
//   - short: Short description of the command, used when printing
//     usage when the command is part of a Menu or when long is empty.
//   - long: Long description of the command, used when printing usage.
//   - subcmds: The subcommands you can invoke through this menu.
func NewMenu(name, short, long string, subcmds ...*Command) *Command {
	cmd, err := NewMenuError(name, short, long, subcmds...)
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}

	return cmd
}

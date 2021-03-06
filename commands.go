package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
	"github.com/mailund/cli/internal/flags"
	"github.com/mailund/cli/internal/params"
	"github.com/mailund/cli/internal/vals"
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

	flags  *flags.FlagSet
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
	for _, sub := range cmd.Subcommands {
		sub.SetOutput(out)
	}

	cmd.out = out
}

// Usage prints usage information about a command.
func (cmd *Command) Usage() { cmd.CommandSpec.Usage() }

// SetUsage sets the function for printing usage information.
// If you only set the Usage field in a spec, it only applies to
// that spec, not the flags or parameters. You
// should almost always use SetUsage() instead.
func (cmd *Command) SetUsage(usage func()) {
	cmd.CommandSpec.Usage = usage
}

// RunError parses options and arguments from args and then executes the
// command. This function returns an error if there are errors parsing or preparing
// the command line arguments, and the command's error handling flag is
// failure.ContinueOnError. You most likely want to use the Run() method
// instead, unless you have good reasons to capture errors rather than
// terminate your program on parsing errors.
func (cmd *Command) RunError(args []string) error {
	if err := cmd.flags.Parse(args); err != nil {
		return err
	}

	if err := cmd.params.Parse(cmd.flags.Args()); err != nil {
		return err
	}

	if err := prepareFlagsAndParams(cmd); err != nil {
		return err
	}

	// Invoke the action for this (sub)command
	if cmd.Action != nil {
		cmd.Action(cmd.argv)
	}

	// then, if there are sub-commands, dispatch
	if len(cmd.subcommands) > 0 {
		subcmd, ok := cmd.subcommands[cmd.command]
		if !ok {
			return interfaces.ParseErrorf("'%s' is not a valid command for %s.\n\n", cmd.command, cmd.Name)
		}

		return subcmd.RunError(cmd.cmdArgs)
	}

	return nil
}

// Run parses options and arguments from args and then executes the
// command.
//
// Parsing errors for either flags or parameters will terminate the program
// with os.Exit(0) for -help options and os.Exit(2) otherwise. If the parsing
// is succesfull, the underlying run callback is executed.
func (cmd *Command) Run(args []string) {
	if err := cmd.RunError(args); err != nil {
		fmt.Fprintf(cmd.out, "Error: %s.\n", err)
		failure.Failure()
	}
}

func showHelp(usage func()) func() error {
	return func() error {
		usage()
		failure.Success()

		return nil
	}
}

// NewCommandError Create a new command. The function returns a new command object or an error.
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
		flags:       flags.NewFlagSet(),
		params:      params.NewParamSet()}

	const linewidth = 70
	cmd.Long = wordWrap(cmd.Long, linewidth)

	if spec.Usage == nil {
		cmd.SetUsage(DefaultUsage(cmd))
	}

	// There is always a help command when we parse, but the usage won't
	// show it unless we make it explicit
	hf := vals.FuncNoValue(showHelp(cmd.Usage))
	_ = cmd.flags.Var(hf, "help", "h", fmt.Sprintf("show help for %s", cmd.Name)) // cannot fail

	if spec.Init != nil {
		cmd.argv = spec.Init()

		if err := connectSpecsFlagsAndParams(cmd, cmd.argv); err != nil {
			return nil, err
		}
	}

	cmd.SetOutput(os.Stdout)

	if len(cmd.Subcommands) > 0 {
		cmd.subcommands = map[string]*Command{}

		for _, sub := range cmd.Subcommands {
			cmd.subcommands[sub.Name] = sub
		}

		cmd.params.Var((*vals.StringValue)(&cmd.command), "cmd", "sub-command to call")
		cmd.params.VariadicVar((*vals.VariadicStringValue)(&cmd.cmdArgs), "...", "argument for sub-commands", 0)
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
			"Usage: %s [flags] %s\n\n",
			cmd.Name, cmd.params.ShortUsage())

		if cmd.Long != "" {
			fmt.Fprintf(cmd.Output(), "%s\n\n", cmd.Long)
		} else {
			fmt.Fprintf(cmd.Output(), "%s\n\n", cmd.Short)
		}

		// Print options and arguments at the bottom.
		fmt.Fprintf(cmd.out, "\n")
		cmd.flags.PrintDefaults(cmd.out)
		fmt.Fprintf(cmd.out, "\n")
		cmd.params.PrintDefaults(cmd.out)

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
	// A menu cannot fail. The subcommands are already parsed, and we don't
	// provide any specs that can fail through the documentation strings.
	cmd, _ := NewCommandError(CommandSpec{
		Name:        name,
		Short:       short,
		Long:        long,
		Subcommands: subcmds,
	})

	return cmd
}

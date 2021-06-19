package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mailund/cli/params"
)

type (
	UsageFunc     = func()
	RunFunc       = func(args []string)
	InitFunc      = func(f *flag.FlagSet, a *params.ParamSet) RunFunc
	InitUsageFunc = func(f *flag.FlagSet, a *params.ParamSet) (UsageFunc, RunFunc)
)

type Command struct {
	Name             string
	ShortDescription string
	LongDescription  string
	Flags            *flag.FlagSet
	Params           *params.ParamSet
	UsageCB          UsageFunc
	RunCB            RunFunc
}

func (cmd Command) Usage() { cmd.UsageCB() }
func (cmd Command) Run(args []string) {
	cmd.Flags.Parse(args)
	cmd.Params.Parse(cmd.Flags.Args())
	cmd.RunCB(cmd.Params.Args())
}

func NewCommandWithUsage(
	name, short_description, long_description string,
	init InitUsageFunc) *Command {
	cmd := &Command{Name: name,
		ShortDescription: short_description,
		LongDescription:  long_description}
	cmd.Flags = flag.NewFlagSet(name, flag.ExitOnError)
	cmd.Params = params.NewParamSet(name) // flags?
	cmd.UsageCB, cmd.RunCB = init(cmd.Flags, cmd.Params)
	cmd.Flags.Usage = cmd.UsageCB
	cmd.Params.Usage = cmd.UsageCB

	// There is always a help command when we parse, but the usage won't
	// show it unless we make it explicit
	cmd.Flags.Bool("help", false, fmt.Sprintf("show help for %s", name))

	return cmd
}

func NewCommand(name, short_description, long_description string, init InitFunc) *Command {
	return NewCommandWithUsage(name, short_description, long_description,
		func(f *flag.FlagSet, p *params.ParamSet) (UsageFunc, RunFunc) {
			return DefaultUsage(name, short_description, long_description, f, p),
				init(f, p)
		})
}

func DefaultUsage(name, description, long_description string,
	f *flag.FlagSet, p *params.ParamSet) UsageFunc {
	return func() {
		fmt.Fprintf(f.Output(),
			"Usage: %s [options] %s\n\n",
			name, strings.Join(p.ParamNames(), " "))
		fmt.Fprintf(f.Output(), "%s\n", description)
		fmt.Fprintf(f.Output(), "Options:\n")
		f.PrintDefaults()
		fmt.Fprintf(f.Output(), "\n")
		p.PrintDefaults()
		if len(long_description) > 0 {
			fmt.Fprintf(f.Output(), "\n%s\n", long_description)
		}
	}
}

func NewMenu(name, short_description, long_description string, subcmds ...*Command) *Command {
	subcommands := map[string]*Command{}
	for _, cmd := range subcmds {
		subcommands[cmd.Name] = cmd
	}

	init := func(f *flag.FlagSet, p *params.ParamSet) (UsageFunc, RunFunc) {
		var command string
		p.StringVar(&command, "cmd", "command to run")

		usage := func() {
			DefaultUsage(name, short_description, "", f, p)()
			fmt.Fprintf(f.Output(), "\nCommands:\n")
			for name, cmd := range subcommands {
				fmt.Fprintf(f.Output(), "\t%s\t%s\n", name, cmd.ShortDescription)
			}
			fmt.Fprintf(f.Output(), "\n")
			if len(long_description) > 0 {
				fmt.Fprintf(f.Output(), "\n%s\n", long_description)
			}
		}
		run := func(args []string) {
			if subcmd, ok := subcommands[command]; ok {
				subcmd.Run(args)
			} else {
				fmt.Fprintf(f.Output(),
					"'%s' is not a valid command for %s.\n\n", command, name)
				f.Usage()
				os.Exit(1)
			}
		}
		return usage, run
	}

	return NewCommandWithUsage(name, short_description, long_description, init)
}

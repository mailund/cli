package flags

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
)

// Flag is the data associated with a single flag.
type Flag struct {
	Name     string               // Name is the short name used for a parameter
	Desc     string               // Desc is a short description of the parameter
	Value    interfaces.FlagValue // Encapsulated value
	DefValue string               // Default value (as string)
}

// FlagSet wraps a set of command line flags.
type FlagSet struct {
	Usage func() // Function for printing usage. Can be changed.
	Name  string // Name of the flag set, usually the same as the command

	flagsList     []*Flag
	flagsMap      map[string]*Flag
	args          []string // arguments after flags
	errorHandling failure.ErrorHandling
	output        io.Writer
}

// SetErrFlag sets the error handling value
func (f *FlagSet) SetErrFlag(eh failure.ErrorHandling) {
	f.errorHandling = eh
}

// SetOutput changes the output stream for the flag set.
func (f *FlagSet) SetOutput(w io.Writer) {
	f.output = w
}

// Lookup gets a flag by name
func (f *FlagSet) Lookup(name string) *Flag {
	return f.flagsMap[name]
}

// Output returns the output stream for the flag set.
func (f *FlagSet) Output() io.Writer { return f.output }

// Args returns the remaining arguments after flags are parsed.
func (f *FlagSet) Args() []string { return f.args }

// NewFlagSet creates a new flag set.
func NewFlagSet(name string, errHandling failure.ErrorHandling) *FlagSet {
	f := &FlagSet{
		Name:          name,
		errorHandling: errHandling,
		flagsList:     []*Flag{},
		flagsMap:      map[string]*Flag{},
		output:        os.Stderr,
	}
	f.Usage = f.PrintDefaults

	return f
}

// Var inserts a new flag in the form of a value.
func (f *FlagSet) Var(value interfaces.FlagValue, name, descr string) error {
	if _, found := f.flagsMap[name]; found {
		return interfaces.SpecErrorf("flag %s is defined more than once", name)
	}

	// Remember the default value as a string; it won't change.
	flag := &Flag{name, descr, value, value.String()}

	f.flagsMap[name] = flag
	f.flagsList = append(f.flagsList, flag)

	return nil
}

// NFlags returns the number of flags in the set.
func (f *FlagSet) NFlags() int {
	return len(f.flagsList)
}

// Flag returns flag number i.
func (f *FlagSet) Flag(i int) *Flag {
	return f.flagsList[i]
}

func (f *FlagSet) sortFlags() {
	sort.Slice(f.flagsList, func(i, j int) bool {
		return f.flagsList[i].Name < f.flagsList[j].Name
	})
}

// PrintDefaults print the default usage for the flags.
func (f *FlagSet) PrintDefaults() {
	if f.NFlags() == 0 {
		return // nothing to print...
	}

	f.sortFlags()
	fmt.Fprintf(f.output, "Flags:\n")

	for _, flag := range f.flagsList {
		defVal := flag.DefValue
		if defVal != "" {
			defVal = " (default " + defVal + ")"
		}

		// FIXME: there is a lot more to do here with optional values and short
		// flags and such...
		fmt.Fprintf(f.output, "  --%s value\n\t%s%s\n", flag.Name, flag.Desc, defVal)
	}
}

func (f *FlagSet) parseShort() error {
	// FIXME
	panic("should not be called yet!")
}

func (f *FlagSet) parseLong() error {
	name := f.args[0][2:]

	if name == "" || name[0] == '-' || name[0] == '=' {
		return interfaces.ParseErrorf("bad flag syntax: %s", f.args[0])
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""

	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]

			break
		}
	}

	flag, alreadythere := f.flagsMap[name]
	if !alreadythere {
		return interfaces.ParseErrorf("flag provided but not defined: --%s", name)
	}

	// FIXME: similar for options that do not take *any* value

	if fv, ok := flag.Value.(interfaces.BoolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := flag.Value.Set(value); err != nil {
				return interfaces.ParseErrorf("error parsing flag %s: %s", name, err)
			}
		} else {
			if err := flag.Value.Set("true"); err != nil {
				return interfaces.ParseErrorf("error parsing flag %s: %s", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return interfaces.ParseErrorf("flag needs an argument: --%s", name)
		}
		if err := flag.Value.Set(value); err != nil {
			return interfaces.ParseErrorf("invalid value %q for flag --%s: %v", value, name, err)
		}
	}

	return nil // success
}

// FIXME: I nicked this from flag, but I need to reimplement it
func (f *FlagSet) parseOne() (more bool, err error) {
	if len(f.args) == 0 {
		return false, nil
	}

	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}

	if s == "--" { // "--" terminates the flags
		f.args = f.args[1:]
		return false, nil
	}

	if strings.HasPrefix(s, "--") {
		if err := f.parseLong(); err != nil {
			return false, err
		}

		return true, nil
	}

	// if we get here, it is a flag, but not a long one...
	if err := f.parseShort(); err != nil {
		return false, err
	}

	return true, nil
}

// Parse parses the flags in the args, leaving the remaining arguments
// in f.Args().
func (f *FlagSet) Parse(args []string) error {
	f.args = args

	for {
		more, err := f.parseOne()

		if err != nil {
			fmt.Fprintf(f.Output(), "Error parsing flag:  %s.", err)

			if f.errorHandling == failure.ExitOnError {
				failure.Failure()
			}

			return err
		}

		if !more {
			return nil
		}
	}
}

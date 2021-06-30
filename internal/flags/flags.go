package flags

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
)

// Flag is the data associated with a single flag.
type Flag struct {
	Long     string               // Name is the long flag name
	Short    string               // Short is the short flag
	Desc     string               // Desc is a short description of the parameter
	Value    interfaces.FlagValue // Encapsulated value
	DefValue string               // Default value (as string)
}

// FlagSet wraps a set of command line flags.
type FlagSet struct {
	Usage func() // Function for printing usage. Can be changed.
	Name  string // Name of the flag set, usually the same as the command

	flagsList []*Flag
	longMap   map[string]*Flag
	shortMap  map[string]*Flag

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
	if f, ok := f.longMap[name]; ok {
		return f
	}

	return f.shortMap[name]
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
		shortMap:      map[string]*Flag{},
		longMap:       map[string]*Flag{},
		output:        os.Stderr,
	}
	f.Usage = f.PrintDefaults

	return f
}

// Var inserts a new flag in the form of a value.
func (f *FlagSet) Var(value interfaces.FlagValue, long, short, descr string) error {
	if len(short) > 1 {
		return interfaces.SpecErrorf("short arguments can only be one character: %s", short)
	}

	if _, found := f.longMap[long]; found {
		return interfaces.SpecErrorf("flag %s is defined more than once", long)
	}

	if _, found := f.shortMap[short]; found {
		return interfaces.SpecErrorf("flag %s is defined more than once", short)
	}

	// Remember the default value as a string; it won't change.
	flag := &Flag{
		Long: long, Short: short, Desc: descr,
		Value: value, DefValue: value.String()}

	if long != "" {
		f.longMap[long] = flag
	}

	if short != "" {
		f.shortMap[short] = flag
	}

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

// PrintDefaults print the default usage for the flags.
func (f *FlagSet) PrintDefaults() {
	if f.NFlags() == 0 {
		return // nothing to print...
	}

	fmt.Fprintf(f.output, "Flags:\n")

	for _, flag := range f.flagsList {
		defVal := flag.DefValue
		if defVal != "" {
			defVal = " (default " + defVal + ")"
		}

		// FIXME: there is a lot more to do here with optional values and short
		// flags and such...
		shortFlag, longFlag := "", ""

		if flag.Short != "" {
			shortFlag = "-" + flag.Short
		}

		if flag.Long != "" {
			if shortFlag != "" {
				shortFlag += ","
			}

			longFlag = "--" + flag.Long
		}

		value := " value"

		if flag.noValues() {
			value = ""
		}

		if def, ok := flag.hasDefault(); ok {
			value = " [value] (default " + def + ")"
		}

		fmt.Fprintf(f.output, "  %s%s%s\n\t%s%s\n", shortFlag, longFlag, value, flag.Desc, defVal)
	}
}

func (f *Flag) noValues() bool {
	if nv, ok := f.Value.(interfaces.NoValueFlag); ok {
		return nv.NoValueFlag()
	}

	return false
}

func (f *Flag) hasDefault() (string, bool) {
	if def, ok := f.Value.(interfaces.DefaultValueFlag); ok {
		return def.Default(), ok
	}

	return "", false
}

func wrapShortParseError(name string, err error) error {
	if err != nil {
		return interfaces.ParseErrorf("error parsing flag -%s: %s", name, err)
	}

	return nil
}

func (f *FlagSet) parseShort() error {
	flags := f.args[0][1:]
	f.args = f.args[1:]

	// Just to avoid some common error...
	for i := 0; i < len(flags); i++ {
		if flags[i] == '=' {
			return interfaces.ParseErrorf("--flag=value syntax for flags only allowed for long options: -%s", flags)
		}
	}

	// every flag except the last cannot a value, so we treat those first
	for i := 0; i < len(flags)-1; i++ {
		x := string(flags[i])
		flag, valid := f.shortMap[x]

		if !valid {
			return interfaces.ParseErrorf("flag provided but not defined: -%s", x)
		} else if flag.noValues() {
			if err := flag.Value.Set(""); err != nil {
				return interfaces.ParseErrorf("error evaluating flag -%s: %s", x, err)
			}
		} else if def, ok := flag.hasDefault(); ok {
			if err := flag.Value.Set(def); err != nil {
				return interfaces.ParseErrorf("error evaluating flag -%s: %s", x, err)
			}
		} else {
			// only the last flag in flags get a value, and this isn't it
			return interfaces.ParseErrorf("flag -%s needs an argument", x)
		}
	}

	x := string(flags[len(flags)-1])
	flag, valid := f.shortMap[x]

	if !valid {
		return interfaces.ParseErrorf("flag provided but not defined: -%s", x)
	}

	if flag.noValues() {
		return wrapShortParseError(x, flag.Value.Set(""))
	}

	if def, ok := flag.hasDefault(); ok {
		return wrapShortParseError(x, flag.Value.Set(def))
	}

	if len(f.args) == 0 || f.args[0][0] == '-' {
		return interfaces.ParseErrorf("flag -%s needs an argument", x)
	}

	// get the next argument as the value for the flag
	value := f.args[0]
	f.args = f.args[1:]

	return wrapShortParseError(x, flag.Value.Set(value))
}

func wrapLongParseError(name string, err error) error {
	if err != nil {
		return interfaces.ParseErrorf("error parsing flag --%s: %s", name, err)
	}

	return nil
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

	flag, valid := f.longMap[name]
	if !valid {
		return interfaces.ParseErrorf("flag provided but not defined: --%s", name)
	}

	if hasValue {
		if flag.noValues() {
			return interfaces.ParseErrorf("flag --%s cannot take values", name)
		}

		return wrapLongParseError(name, flag.Value.Set(value))
	}

	if flag.noValues() {
		// we don't take values, so we can stop with this flag. Invoke it by
		// calling Set() with the empty string
		return wrapLongParseError(name, flag.Value.Set(""))
	}

	if def, ok := flag.hasDefault(); ok {
		// flags with defaults can only take arguments as --flag=arg,
		// since we cannot figure out if --flag xxx means that xxx is a
		// flag argument or a positional argument. So if we have one of
		// those, then we invoke it here, and do not look at the following
		// arg.
		return wrapLongParseError(name, flag.Value.Set(def))
	}

	if len(f.args) == 0 || f.args[0][0] == '-' {
		return interfaces.ParseErrorf("flag --%s needs an argument", name)
	}

	// get the next argument as the value for the flag
	value, f.args = f.args[0], f.args[1:]

	return wrapLongParseError(name, flag.Value.Set(value))
}

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

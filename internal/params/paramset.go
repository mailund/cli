package params

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mailund/cli/inter"
	"github.com/mailund/cli/internal/failure"
)

type (
	/* ErrorHandling flags for determining how the argument
	parser should handle errors.

	The error handling works the same as for the flag package,
	except that ParamSet's only support two flags:

		- ExitOnError: Terminate the program on error.
		- ContinueOnError: Report the error as an error object.
	*/
	ErrorHandling = flag.ErrorHandling
)

const (
	// ContinueOnError means that parsing will return an error
	// rather than terminate the program.
	ContinueOnError = flag.ContinueOnError
	// ExitOnError means that parsing will exit the program
	// (using os.Exit(2)) if the parser fails.
	ExitOnError = flag.ExitOnError
	// PanicOnError means we raise a panic on errors
	PanicOnError = flag.PanicOnError
)

// Param holds positional parameter information.
type Param struct {
	// Name is the short name used for a parameter
	Name string
	// Desc is a short description of the parameter
	Desc string
	// Encapsulated value
	Value inter.PosValue
}

// VariadicParam holds information about a variadic argument.
type VariadicParam struct {
	// Name is the short name used for a parameter
	Name string
	// Desc is a short description of the parameter
	Desc string
	// Min is the minimum number of parameters that the variadic
	// parameter takes.
	Min int

	parser func([]string) error
}

// ParamSet contains a list of specified parameters for a
// commandline tool and parsers for parsing commandline
// arguments.
type ParamSet struct {
	// Name is a name used when printing usage information
	// for a parameter set.
	Name string
	// Usage is a function you can assign to for changing
	// the help info.
	Usage func()
	// ErrorFlag Controls how we deal with parsing errors
	ErrorFlag ErrorHandling

	params []*Param
	out    io.Writer

	// last parameter, used for variadic arguments
	last *VariadicParam
}

// SetFlag sets the error handling flag.
func (p *ParamSet) SetFlag(f ErrorHandling) {
	p.ErrorFlag = f
}

// SetOutput specifies where usage output and error messages should
// be written to.
//
// Parameters:
//   - out: a writer object. The default, if you do not set a new
//     variable is os.Stderr.
func (p *ParamSet) SetOutput(out io.Writer) { p.out = out }

// Output returns the output stream where messages are written to
func (p *ParamSet) Output() io.Writer { return p.out }

// PrintDefaults prints a description of the parameters
func (p *ParamSet) PrintDefaults() {
	if p.NParams() == 0 && p.last == nil {
		return // nothing to print...
	}

	fmt.Fprintf(p.Output(), "Arguments:\n")

	for _, par := range p.params {
		fmt.Fprintf(p.Output(), "  %s\n\t%s\n", par.Name, par.Desc)
	}

	if p.last != nil {
		fmt.Fprintf(p.Output(), "  %s\n\t%s\n", p.last.Name, p.last.Desc)
	}
}

// NParams returns the number of parameters in the set, excluding the last variadic
// parameter if there is one. You can test for whether there is a variadic argument using
// p.Variadic() != nil.
func (p *ParamSet) NParams() int {
	return len(p.params)
}

// Param returns the parameter at position i.
func (p *ParamSet) Param(i int) *Param {
	return p.params[i]
}

// Variadic returns the last variadic parameter, if there is one.
func (p *ParamSet) Variadic() *VariadicParam {
	return p.last
}

// NewParamSet creates a new parameter set.
//
// Parameters:
//   - name: The name for the parameter set. Used when printing
//     usage and error information.
//   - errflat: Controls the error handling. If ExitOnError, parse
//     errors will terminate the program (os.Exit(1)); if ContinueOnError
//     the parsing will return an error instead.
func NewParamSet(name string, errflag ErrorHandling) *ParamSet {
	argset := &ParamSet{
		Name:      name,
		ErrorFlag: errflag,
		params:    []*Param{},
		out:       os.Stderr}
	argset.Usage = argset.PrintDefaults

	return argset
}

// ShortUsage returns a string used for printing the usage
// of a parameter set.
func (p *ParamSet) ShortUsage() string {
	names := make([]string, len(p.params))
	for i, param := range p.params {
		names[i] = param.Name
	}

	namesUsage := strings.Join(names, " ")

	if p.last != nil {
		namesUsage += " " + p.last.Name
	}

	return namesUsage
}

// Parse parses arguments against parameters.
//
// Parameters:
//   - args: A slice of strings from a command line.
//
// If the parameter set has error handling flag ExitOnError,
// a parsing error will terminate the program (os.Exit(1)).
// If the parameter set has error handling flag ContinueOnError,
// it will return an error instead. If all goes well, it will
// return nil.
func (p *ParamSet) Parse(args []string) error {
	minParams := len(p.params)
	if p.last != nil {
		minParams += p.last.Min
	}

	if len(args) < minParams {
		if p.ErrorFlag == ExitOnError {
			fmt.Fprintf(p.Output(),
				"Too few arguments for command '%s'\n\n",
				p.Name)
			p.Usage()
			failure.Failure()
		}

		return inter.ParseErrorf("too few arguments")
	}

	if p.last == nil && len(args) > len(p.params) {
		if p.ErrorFlag == ExitOnError {
			fmt.Fprintf(p.Output(),
				"Too many arguments for command '%s'\n\n",
				p.Name)
			p.Usage()
			failure.Failure()
		}

		return inter.ParseErrorf("too many arguments")
	}

	for i, par := range p.params {
		if err := par.Value.Set(args[i]); err != nil {
			if p.ErrorFlag == flag.ExitOnError {
				fmt.Fprintf(p.Output(), "Error parsing parameter %s='%s', %s.\n",
					par.Name, args[i], err.Error())
				failure.Failure()
			}

			return inter.ParseErrorf("error parsing parameter %s='%s'", par.Name, args[i])
		}
	}

	if p.last != nil {
		rest := args[len(p.params):]
		if err := p.last.parser(rest); err != nil {
			if p.ErrorFlag == ExitOnError {
				fmt.Fprintf(p.Output(), "Error parsing parameters %s='%v', %s.\n",
					p.last.Name, rest, err.Error())
				failure.Failure()
			}

			return inter.ParseErrorf("error parsing parameters %s='%v'", p.last.Name, rest)
		}
	}

	return nil
}

// Var adds a new PosValue variable to the parameter set.
//
// Parameters:
//   - val: a variable where the parsed argument should be written.
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
func (p *ParamSet) Var(val inter.PosValue, name, desc string) {
	p.params = append(p.params, &Param{name, desc, val})
}

func variadicStringParser(target *[]string) func([]string) error {
	return func(args []string) error {
		*target = args
		return nil
	}
}

func variadicBoolParser(target *[]bool) func([]string) error {
	return func(args []string) error {
		res := make([]bool, len(args))

		for i, x := range args {
			val, err := strconv.ParseBool(x)
			if err != nil {
				return inter.ParseErrorf("cannot parse '%s' as boolean", x)
			}

			res[i] = val
		}

		*target = res

		return nil
	}
}

func variadicIntParser(target *[]int) func([]string) error {
	return func(args []string) error {
		res := make([]int, len(args))

		for i, x := range args {
			val, err := strconv.ParseInt(x, 0, 64) //nolint:gomnd // 64 bits is not magic
			if err != nil {
				return inter.ParseErrorf("cannot parse '%s' as integer", x)
			}

			res[i] = int(val)
		}

		*target = res

		return nil
	}
}

func variadicFloatParser(target *[]float64) func([]string) error {
	return func(args []string) error {
		res := make([]float64, len(args))

		for i, x := range args {
			val, err := strconv.ParseFloat(x, 64) //nolint:gomnd // 64 bits is not magic
			if err != nil {
				return inter.ParseErrorf("cannot parse '%s' as float", x)
			}

			res[i] = val
		}

		*target = res

		return nil
	}
}

// VariadicStringVar install a variadic string argument
// as the last parameter(s) for the parameter set.
//
// Parameters:
//   - target: Pointer to where the parsed arguments should be written.
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicStringVar(target *[]string, name, desc string, min int) {
	p.VariadicFunc(name, desc, min, variadicStringParser(target))
}

// VariadicString install a variadic string argument
// as the last parameter(s) for the parameter set. It returns a pointer
// to where the parsed values will go if parsing is succesfull.
//
// Parameters:
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicString(name, desc string, min int) *[]string {
	var x = []string{}

	p.VariadicStringVar(&x, name, desc, min)

	return &x
}

// VariadicBoolVar install a variadic bool argument
// as the last parameter(s) for the parameter set.
//
// Parameters:
//   - target: Pointer to where the parsed arguments should be written.
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicBoolVar(target *[]bool, name, desc string, min int) {
	p.VariadicFunc(name, desc, min, variadicBoolParser(target))
}

// VariadicBool install a variadic bool argument
// as the last parameter(s) for the parameter set. It returns a pointer
// to where the parsed values will go if parsing is succesfull.
//
// Parameters:
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicBool(name, desc string, min int) *[]bool {
	var x = []bool{}

	p.VariadicBoolVar(&x, name, desc, min)

	return &x
}

// VariadicIntVar install a variadic int argument
// as the last parameter(s) for the parameter set.
//
// Parameters:
//   - target: Pointer to where the parsed arguments should be written.
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicIntVar(target *[]int, name, desc string, min int) {
	p.VariadicFunc(name, desc, min, variadicIntParser(target))
}

// VariadicInt install a variadic int argument
// as the last parameter(s) for the parameter set. It returns a pointer
// to where the parsed values will go if parsing is succesfull.
//
// Parameters:
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicInt(name, desc string, min int) *[]int {
	var x = []int{}

	p.VariadicIntVar(&x, name, desc, min)

	return &x
}

// VariadicFloatVar install a variadic float argument
// as the last parameter(s) for the parameter set.
//
// Parameters:
//   - target: Pointer to where the parsed arguments should be written.
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicFloatVar(target *[]float64, name, desc string, min int) {
	p.VariadicFunc(name, desc, min, variadicFloatParser(target))
}

// VariadicFloat install a variadic float argument
// as the last parameter(s) for the parameter set. It returns a pointer
// to where the parsed values will go if parsing is succesfull.
//
// Parameters:
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicFloat(name, desc string, min int) *[]float64 {
	var x = []float64{}

	p.VariadicFloatVar(&x, name, desc, min)

	return &x
}

// VariadicFunc install a variadic callback function
// as the last parameter(s) for the parameter set. The callback will be
// invoked when the parser reaches the end of the normal parameters.
//
// Parameters:
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
//   - fn: Callback function invoked on the last arguments. If parsing is
//     succesfull it should return nil, otherwise a non-nil error.
func (p *ParamSet) VariadicFunc(name, desc string, min int, fn func([]string) error) {
	p.last = &VariadicParam{Name: name, Desc: desc, Min: min, parser: fn}
}

package params

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
)

func paramDescription(val interface{}, descr string) string {
	if d, ok := val.(interfaces.ArgumentDescription); ok {
		return d.ArgumentDescription(false, descr)
	}

	return descr
}

// Param holds positional parameter information.
type Param struct {
	// Name is the short name used for a parameter
	Name string
	// Desc is a short description of the parameter
	Desc string
	// Encapsulated value
	Value interfaces.PosValue
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
	// Encapsulated value
	Value interfaces.VariadicValue
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
	ErrorFlag failure.ErrorHandling

	params []*Param
	out    io.Writer

	// last parameter, used for variadic arguments
	last *VariadicParam
}

// SetErrFlag sets the error handling flag.
func (p *ParamSet) SetErrFlag(f failure.ErrorHandling) {
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
func NewParamSet(name string, errflag failure.ErrorHandling) *ParamSet {
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
		if p.ErrorFlag == failure.ExitOnError {
			fmt.Fprintf(p.Output(),
				"Too few arguments for command '%s'\n\n",
				p.Name)
			p.Usage()
			failure.Failure()
		}

		return interfaces.ParseErrorf("too few arguments")
	}

	if p.last == nil && len(args) > len(p.params) {
		if p.ErrorFlag == failure.ExitOnError {
			fmt.Fprintf(p.Output(),
				"Too many arguments for command '%s'\n\n",
				p.Name)
			p.Usage()
			failure.Failure()
		}

		return interfaces.ParseErrorf("too many arguments")
	}

	for i, par := range p.params {
		if err := par.Value.Set(args[i]); err != nil {
			if p.ErrorFlag == failure.ExitOnError {
				fmt.Fprintf(p.Output(), "Error parsing parameter %s='%s', %s.\n",
					par.Name, args[i], err.Error())
				failure.Failure()
			}

			return interfaces.ParseErrorf("error parsing parameter %s='%s'", par.Name, args[i])
		}
	}

	if p.last != nil {
		rest := args[len(p.params):]
		if err := p.last.Value.Set(rest); err != nil {
			if p.ErrorFlag == failure.ExitOnError {
				fmt.Fprintf(p.Output(), "Error parsing parameters %s='%v', %s.\n",
					p.last.Name, rest, err.Error())
				failure.Failure()
			}

			return interfaces.ParseErrorf("error parsing parameters %s='%v'", p.last.Name, rest)
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
func (p *ParamSet) Var(val interfaces.PosValue, name, desc string) {
	p.params = append(p.params, &Param{name, paramDescription(val, desc), val})
}

// VariadicVar install a variadic argument
// as the last parameter(s) for the parameter set.
//
// Parameters:
//   - val: A variable that will hold the parsed arguments.
//   - name: Name of the argument, used when printing usage.
//   - desc: Description of the argument. Used when printing usage.
//   - min: The minimum number of arguments that the command line must
//     have for this parameter.
func (p *ParamSet) VariadicVar(val interfaces.VariadicValue, name, desc string, min int) {
	p.last = &VariadicParam{Name: name, Desc: paramDescription(val, desc), Min: min, Value: val}
}

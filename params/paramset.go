package params

import (
	"fmt"
	"io"
	"os"
)

type (
	UsageFunc = func()
	Parser    = func(arg string)
)

type Param struct {
	name   string
	desc   string
	parser Parser
}

// FIXME: find a way to specify whether you can have more of an
// argument or whether the argset must match completely

type ParamSet struct {
	name   string
	params []*Param
	args   []string // will be set to remaining args after parsing
	out    io.Writer

	// You can assign to usage to change the help info.
	Usage UsageFunc

	// String that comes after the parameter names in usage strings
	// if you set it
	Suffix string
}

func (p *ParamSet) Output() io.Writer { return p.out }

func (p *ParamSet) PrintDefaults() {
	fmt.Fprintf(p.out, "Arguments:\n")
	for _, param := range p.params {
		fmt.Fprintf(p.out, "  %s\n\t%s\n", param.name, param.desc)
	}
}

func NewParamSet(name string) *ParamSet {
	argset := &ParamSet{name: name,
		params: []*Param{}, args: []string{}}
	argset.Usage = argset.PrintDefaults
	argset.out = os.Stderr
	return argset
}

func (p *ParamSet) ParamNames() []string {
	n := len(p.params)
	if p.Suffix != "" {
		n++
	}
	names := make([]string, n)
	for i, param := range p.params {
		names[i] = param.name
	}
	if p.Suffix != "" {
		names[n-1] = p.Suffix
	}
	return names
}

func (p *ParamSet) Parse(args []string) {
	if len(args) < len(p.params) {
		fmt.Fprintf(p.out,
			"Insufficient arguments for command '%s'\n\n",
			p.name)
		p.Usage()
		os.Exit(1) // FIXME: control by flag?
	}
	for i, param := range p.params {
		param.parser(args[i])
	}
	p.args = args[len(p.params):]
}

func (p *ParamSet) Args() []string { return p.args }

func stringParser(target *string) Parser {
	return func(arg string) { *target = arg }
}

func (p *ParamSet) StringVar(target *string, name, desc string) {
	param := &Param{name, desc, stringParser(target)}
	p.params = append(p.params, param)
}

func (p *ParamSet) SetSuffix(suffix string) {
	p.Suffix = suffix
}

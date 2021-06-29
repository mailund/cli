package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type typeSpec struct {
	TypeName string
	Parse    string
	Format   string

	SetInput     string
	SetOutput    string
	CantFail     bool
	SetFailInput string

	VarInput     string
	VarOutput    string
	VarFailInput string
}

var valTypes = []typeSpec{
	{
		TypeName:  "string",
		Parse:     "x",
		Format:    "string(*val)",
		SetInput:  "foo",
		SetOutput: "foo",
		CantFail:  true,
		VarInput:  `"foo", "bar", "baz"`,
		VarOutput: `"foo", "bar", "baz"`,
	},
	{
		TypeName:     "bool",
		Parse:        "strconv.ParseBool(x)",
		Format:       "strconv.FormatBool(bool(*val))",
		SetInput:     "true",
		SetOutput:    "true",
		SetFailInput: "foo",
		VarInput:     `"true", "false", "true"`,
		VarOutput:    `true, false, true`,
		VarFailInput: `"foo"`,
	},

	{
		TypeName:     "int",
		Parse:        "strconv.ParseInt(x, 0, strconv.IntSize)",
		Format:       "strconv.Itoa(int(*val))",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "foo",
		VarInput:     `"-1", "2", "-3"`,
		VarOutput:    `-1, 2, -3`,
		VarFailInput: `"foo"`,
	},
	{
		TypeName:     "int8",
		Parse:        "strconv.ParseInt(x, 0, 8)",
		Format:       "strconv.Itoa(int(*val))",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "foo",
		VarInput:     `"-1", "2", "-3"`,
		VarOutput:    `-1, 2, -3`,
		VarFailInput: `"foo"`,
	},
	{
		TypeName:     "int16",
		Parse:        "strconv.ParseInt(x, 0, 16)",
		Format:       "strconv.Itoa(int(*val))",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "foo",
		VarInput:     `"-1", "2", "-3"`,
		VarOutput:    `-1, 2, -3`,
		VarFailInput: `"foo"`,
	},
	{
		TypeName:     "int32",
		Parse:        "strconv.ParseInt(x, 0, 32)",
		Format:       "strconv.Itoa(int(*val))",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "foo",
		VarInput:     `"-1", "2", "-3"`,
		VarOutput:    `-1, 2, -3`,
		VarFailInput: `"foo"`,
	},
	{
		TypeName:     "int64",
		Parse:        "strconv.ParseInt(x, 0, 64)",
		Format:       "strconv.Itoa(int(*val))",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "foo",
		VarInput:     `"-1", "2", "-3"`,
		VarOutput:    `-1, 2, -3`,
		VarFailInput: `"foo"`,
	},

	{
		TypeName:     "uint",
		Parse:        "strconv.ParseUint(x, 0, strconv.IntSize)",
		Format:       "strconv.FormatUint(uint64(*val), 10)",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "-1",
		VarInput:     `"1", "2", "3"`,
		VarOutput:    `1, 2, 3`,
		VarFailInput: `"-1"`,
	},
	{
		TypeName:     "uint8",
		Parse:        "strconv.ParseUint(x, 0, 8)",
		Format:       "strconv.FormatUint(uint64(*val), 10)",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "-1",
		VarInput:     `"1", "2", "3"`,
		VarOutput:    `1, 2, 3`,
		VarFailInput: `"-1"`,
	},
	{
		TypeName:     "uint16",
		Parse:        "strconv.ParseUint(x, 0, 16)",
		Format:       "strconv.FormatUint(uint64(*val), 10)",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "-1",
		VarInput:     `"1", "2", "3"`,
		VarOutput:    `1, 2, 3`,
		VarFailInput: `"-1"`,
	},
	{
		TypeName:     "uint32",
		Parse:        "strconv.ParseUint(x, 0, 32)",
		Format:       "strconv.FormatUint(uint64(*val), 10)",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "-1",
		VarInput:     `"1", "2", "3"`,
		VarOutput:    `1, 2, 3`,
		VarFailInput: `"-1"`,
	},
	{
		TypeName:     "uint64",
		Parse:        "strconv.ParseUint(x, 0, 64)",
		Format:       "strconv.FormatUint(uint64(*val), 10)",
		SetInput:     "42",
		SetOutput:    "42",
		SetFailInput: "-1",
		VarInput:     `"1", "2", "3"`,
		VarOutput:    `1, 2, 3`,
		VarFailInput: `"-1"`,
	},

	{
		TypeName:     "float32",
		Parse:        "strconv.ParseFloat(x, 32)",
		Format:       "strconv.FormatFloat(float64(*val), 'g', -1, 32)",
		SetInput:     "3.14",
		SetOutput:    "3.14",
		SetFailInput: "foo",
		VarInput:     `"0.1", "0.2", "0.3"`,
		VarOutput:    `0.1, 0.2, 0.3`,
		VarFailInput: `"foo"`,
	},
	{
		TypeName:     "float64",
		Parse:        "strconv.ParseFloat(x, 64)",
		Format:       "strconv.FormatFloat(float64(*val), 'g', -1, 64)",
		SetInput:     "3.14",
		SetOutput:    "3.14",
		SetFailInput: "foo",
		VarInput:     `"0.1", "0.2", "0.3"`,
		VarOutput:    `0.1, 0.2, 0.3`,
		VarFailInput: `"foo"`,
	},

	{
		TypeName:     "complex64",
		Parse:        "strconv.ParseComplex(x, 64)",
		Format:       "strconv.FormatComplex(complex128(*val), 'g', -1, 64)",
		SetInput:     "(3.14+42i)",
		SetOutput:    "(3.14+42i)",
		SetFailInput: "foo",
		VarInput:     `"0.1+0.2i", "0.2+0.3i", "0.3+0.4i"`,
		VarOutput:    `0.1+0.2i, 0.2+0.3i, 0.3+0.4i`,
		VarFailInput: `"foo"`,
	},
	{
		TypeName:     "complex128",
		Parse:        "strconv.ParseComplex(x, 128)",
		Format:       "strconv.FormatComplex(complex128(*val), 'g', -1, 128)",
		SetInput:     "(3.14+42i)",
		SetOutput:    "(3.14+42i)",
		SetFailInput: "foo",
		VarInput:     `"0.1+0.2i", "0.2+0.3i", "0.3+0.4i"`,
		VarOutput:    `0.1+0.2i, 0.2+0.3i, 0.3+0.4i`,
		VarFailInput: `"foo"`,
	},
}

var valsTemplate = `
type {{TypeName .TypeName}} {{.TypeName}}

func (val *{{TypeName .TypeName}}) Set(x string) error {
	{{if .CantFail}}*val = {{TypeName .TypeName}}({{.Parse}})
	return nil{{else}}v, err := {{.Parse}}
	if err != nil {
		err = inter.ParseErrorf("argument \"%s\" cannot be parsed as {{.TypeName}}", x)
	} else {
		*val = {{TypeName .TypeName}}(v)
	}

	return err{{end}}
}

func (val *{{TypeName .TypeName}}) String() string {
	return {{.Format}}
}

func {{TypeName .TypeName}}Constructor(val reflect.Value) inter.FlagValue {
	return (*{{TypeName .TypeName}})(val.Interface().(*{{.TypeName}}))
}

type {{VariadicTypeName .TypeName}} []{{.TypeName}}

func (vals *{{VariadicTypeName .TypeName}}) Set(xs []string) error {
	*vals = make([]{{.TypeName}}, len(xs))

	for i, x := range xs {
		{{if .CantFail}}(*vals)[i] = {{.TypeName}}({{.Parse}}){{else}}val, err := {{.Parse}}
		if err != nil {
			return inter.ParseErrorf("cannot parse '%s' as {{.TypeName}}", x)
		}

		(*vals)[i] = {{.TypeName}}(val){{end}}
	}

	return nil
}

func {{VariadicTypeName .TypeName}}Constructor(val reflect.Value) inter.VariadicValue {
	return (*{{VariadicTypeName .TypeName}})(val.Interface().(*[]{{.TypeName}}))
}
`

var tableTemplate = "\tValsConstructors[reflect.TypeOf((*{{.TypeName}})(nil))] = {{TypeName .TypeName}}Constructor\n" +
	"\tVarValsConstructors[reflect.TypeOf((*[]{{.TypeName}})(nil))] = {{VariadicTypeName .TypeName}}Constructor\n"

var testTemplate = `
func Test{{TypeName .TypeName}}(t *testing.T) {
	var (
		x     {{.TypeName}}
		val = vals.AsValue(reflect.ValueOf(&x))
	)

	if val == nil {
		t.Fatal("val should not be nil")
	}

	if err := val.Set("{{.SetInput}}"); err != nil {
		t.Error("error setting val to {{.SetInput}}")
	}

	if val.String() != "{{.SetOutput}}" {
		t.Errorf("Unexpected string value for val: %s", val.String())
	}{{ if (not .CantFail) }}

	if err := val.Set("{{.SetFailInput}}"); err == nil {
		t.Error("val.Set() should fail this time")
	}{{end}}
}

func Test{{VariadicTypeName .TypeName}}(t *testing.T) {
	var (
		x    []{{.TypeName}}
		vv = vals.AsVariadicValue(reflect.ValueOf(&x))
	)

	if vv == nil {
		t.Fatal("vv should not be nil")
	}

	if err := vv.Set([]string{ {{.VarInput}} }); err != nil {
		t.Error("vv.Set should not fail")
	}

	if !reflect.DeepEqual(x, []{{.TypeName}}{ {{.VarOutput}} }) {
		t.Error("x holds the wrong value")
	}{{ if (not .CantFail) }}

	if err := vv.Set([]string{ {{.VarFailInput}} }); err == nil {
		t.Error("vv.Set() should fail this time")
	}{{end}}
}
`

var (
	funcMap = template.FuncMap{
		"TypeName":         func(name string) string { return strings.Title(name) + "Value" },
		"VariadicTypeName": func(name string) string { return "Variadic" + strings.Title(name) + "Value" },
	}

	valsFuncs = template.Must(template.New("val").Funcs(funcMap).Parse(valsTemplate))
	tableInit = template.Must(template.New("val").Funcs(funcMap).Parse(tableTemplate))
	tests     = template.Must(template.New("val").Funcs(funcMap).Parse(testTemplate))
)

func genvals(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fmt.Fprintln(f, `// Code generated by gen/genvals.go DO NOT EDIT.
package vals

import (
	"reflect"
	"strconv"

	"github.com/mailund/cli/inter"
)`)

	for i := 0; i < len(valTypes); i++ {
		_ = valsFuncs.Execute(f, valTypes[i])
	}

	fmt.Fprintln(f, "\nfunc init() {")

	for i := 0; i < len(valTypes); i++ {
		_ = tableInit.Execute(f, valTypes[i])
	}

	fmt.Fprintln(f, "}")
}

func genvalsTest(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fmt.Fprintln(f, `// Code generated by gen/genvals.go DO NOT EDIT.
package vals_test

import (
	"reflect"
	"testing"

	"github.com/mailund/cli/internal/vals"
)`)

	for i := 0; i < len(valTypes); i++ {
		_ = tests.Execute(f, valTypes[i])
	}
}

// I would love to pipe the result through gofmt and goimports,
// but fucking go won't let me install goimports. Stupid stupid stupid go.
func main() {
	genvals(os.Args[1] + ".go")
	genvalsTest(os.Args[1] + "_test.go")
}

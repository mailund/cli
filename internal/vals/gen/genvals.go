package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var valsTemplate = `
type {{TypeName .TypeName}} {{.TypeName}}

func (val *{{TypeName .TypeName}}) Set(x string) error {
	v, err := {{.Parse}}
	if err != nil {
		err = params.ParseErrorf("argument \"%s\" cannot be parsed as {{.TypeName}}", x)
	} else {
		*val = {{TypeName .TypeName}}(v)
	}

	return err
}

func (val *{{TypeName .TypeName}}) String() string {
	return {{.Format}}
}
`

type typeSpec struct {
	TypeName string
	Parse    string
	Format   string
}

var valTypes = []typeSpec{
	{TypeName: "string",
		Parse:  "x, (error)(nil)",
		Format: "`\"` + string(*val) + `\"`"},
	{TypeName: "bool",
		Parse:  "strconv.ParseBool(x)",
		Format: "strconv.FormatBool(bool(*val))"},
	{TypeName: "int",
		Parse:  "strconv.ParseInt(x, 0, strconv.IntSize)",
		Format: "strconv.Itoa(int(*val))"},
	{TypeName: "float64",
		Parse:  "strconv.ParseFloat(x, 64)",
		Format: "strconv.FormatFloat(float64(*val), 'g', -1, 64)"},
}

// I would love to pipe the result through gofmt and goimports,
// but fucking go won't let me install goimports. Stupid stupid stupid go.
func main() {
	t := template.Must(template.New("val").Funcs(template.FuncMap{
		"TypeName": func(name string) string { return strings.Title(name) + "Value" },
	}).Parse(valsTemplate))

	fmt.Println(`// Code generated by genvals.go DO NOT EDIT.
package vals

import (
	"strconv"

	"github.com/mailund/cli/internal/params"
)`)

	for _, tspec := range valTypes {
		_ = t.Execute(os.Stdout, tspec)
	}
}

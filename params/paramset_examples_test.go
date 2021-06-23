package params_test

import (
	"fmt"
	"os"

	"github.com/mailund/cli/params"
)

func ExampleParamSet() {
	p := params.NewParamSet("name", params.ExitOnError)

	// Variables we want to parse
	var (
		i int
		b bool
	)

	p.IntVar(&i, "I", "an integer argument")
	p.BoolVar(&b, "B", "a boolean argument")
	f := p.Float("F", "a float argument")
	s := p.String("S", "a string argument")

	// When we parse the parameter set, i, b, f, and s will
	// get the parameters from the list of arguments.
	_ = p.Parse([]string{"42", "t", "3.14", "foo"})

	fmt.Printf("I=%d, B=%t, F=%f, S='%s'\n", i, b, *f, *s)

	// Output: I=42, B=true, F=3.140000, S='foo'
}

func ExampleParamSet_second() { //nolint:funlen // Example functions can be long as well... (for now)
	p := params.NewParamSet("sum", params.ExitOnError)

	// Variables we want to parse
	var (
		op   string
		args []float64
	)

	p.StringVar(&op, "op", "operation to perform")
	p.VariadicFloatVar(&args, "args", "arg [args...]", 0)

	_ = p.Parse([]string{"+", "0.42", "3.14", "0.3"})

	switch op {
	case "+":
		res := 0.0
		for _, x := range args {
			res += x
		}

		fmt.Printf("Result: %f\n", res)

	case "*":
		res := 1.0
		for _, x := range args {
			res *= x
		}

		fmt.Printf("Result: %f\n", res)

	default:
		fmt.Fprintf(os.Stderr, "I didn't implement any more ops!\n")
		os.Exit(2)
	}

	_ = p.Parse([]string{"*", "0.42", "3.14", "0.3"})

	switch op {
	case "+":
		res := 0.0
		for _, x := range args {
			res += x
		}

		fmt.Printf("Result: %f\n", res)

	case "*":
		res := 1.0
		for _, x := range args {
			res *= x
		}

		fmt.Printf("Result: %f\n", res)

	default:
		fmt.Fprintf(os.Stderr, "I didn't implement any more ops!\n")
		os.Exit(2)
	}

	// Output: Result: 3.860000
	// Result: 0.395640
}

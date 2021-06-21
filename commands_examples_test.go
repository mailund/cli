package cli_test

import (
	"fmt"
	"math"

	"github.com/mailund/cli"
)

func ExampleCommand() {
	cmd := cli.NewCommand(
		"sum", "adds floating point arguments",
		"<long description of how addition works>",
		func(cmd *cli.Command) func() {
			// Arguments to the command
			var round bool
			var args []float64

			// Install arguments, one as a flag and one as a positional arg
			cmd.Flags.BoolVar(&round, "round", false, "round off result")
			cmd.Params.VariadicFloatVar(&args, "args", "arg [args...]", 0)

			// Then return the function for executing the command
			return func() {
				res := 0.0
				for _, x := range args {
					res += x
				}
				if round {
					res = math.Round(res)
				}
				fmt.Printf("Result: %f\n", res)
			}
		})

	// Don't round off the result
	cmd.Run([]string{"0.42", "3.14", "0.3"})

	// Round off the result
	cmd.Run([]string{"-round", "0.42", "3.14", "0.3"})

	// Output:
	// Result: 3.860000
	// Result: 4.000000
}

func ExampleNewMenu() {
	type operator = func(float64, float64) float64
	type initFunc = func(cmd *cli.Command) func()

	apply := func(init float64, op operator, args []float64) float64 {
		res := init
		for _, x := range args {
			res = op(res, x)
		}
		return res
	}
	initCmd := func(init float64, op operator) initFunc {
		return func(cmd *cli.Command) func() {
			var args []float64
			cmd.Params.VariadicFloatVar(&args, "args", "arg [args...]", 0)
			return func() {
				fmt.Printf("Result: %f\n", apply(init, op, args))
			}
		}
	}

	sum := cli.NewCommand(
		"+", "adds floating point arguments",
		"<long description of how addition works>",
		initCmd(0, func(x, y float64) float64 { return x + y }))
	mult := cli.NewCommand(
		"*", "multiplies floating point arguments",
		"<long description of how multiplication works>",
		initCmd(1, func(x, y float64) float64 { return x * y }))

	calc := cli.NewMenu("calc", "does calculations",
		"<long description of how math works>",
		sum, mult)

	// Will call the sum command
	calc.Run([]string{"+", "0.42", "3.14", "0.3"})
	// Will call the mult command
	calc.Run([]string{"*", "0.42", "3.14", "0.3"})

	// Output:
	// Result: 3.860000
	// Result: 0.395640
}

func ExampleNewMenu_second() {
	say := func(x string) func(*cli.Command) func() {
		return func(*cli.Command) func() { return func() { fmt.Println(x) } }
	}
	menu := cli.NewMenu("menu", "", "",
		cli.NewMenu("foo", "", "",
			cli.NewCommand("x", "", "", say("menu/foo/x")),
			cli.NewCommand("y", "", "", say("menu/foo/y"))),
		cli.NewMenu("bar", "", "",
			cli.NewCommand("x", "", "", say("menu/bar/x")),
			cli.NewCommand("y", "", "", say("menu/bar/y"))),
		cli.NewCommand("baz", "", "", say("menu/baz")))

	menu.Run([]string{"baz"})      // will say menu/baz
	menu.Run([]string{"foo", "x"}) // will say menu/foo/x
	menu.Run([]string{"bar", "y"}) // will say menu/bar/y

	// Output:
	// menu/baz
	// menu/foo/x
	// menu/bar/y
}

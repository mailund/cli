package cli_test

import (
	"fmt"
	"math"

	"github.com/mailund/cli"
)

func ExampleCommand() {
	type AddArgs struct {
		X int `arg:"x" descr:"first addition argument"`
		Y int `arg:"y" descr:"first addition argument"`
	}

	add := cli.NewCommand(
		cli.CommandSpec{
			Name:  "add",
			Short: "adds two floating point arguments",
			Long:  "<long description of how addition works>",
			Init:  func() interface{} { return new(AddArgs) },
			Action: func(args interface{}) {
				// Get the argumens through casting the args
				argv, _ := args.(*AddArgs)

				// Then we can do our calculations with the arguments
				fmt.Printf("Result: %d\n", argv.X+argv.Y)
			},
		})

	add.Run([]string{"2", "2"})

	// Output: Result: 4
}

func ExampleCommand_second() {
	type NumArgs struct {
		Round bool      `flag:"round" descr:"should result be rounded to nearest integer"`
		Args  []float64 `arg:"args" descr:"numerical arguments to command"`
	}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:  "sum",
			Short: "adds floating point arguments",
			Long:  "<long description of how addition works>",
			Init:  func() interface{} { return new(NumArgs) },
			Action: func(args interface{}) {
				argv, _ := args.(*NumArgs)
				res := 0.0
				for _, x := range argv.Args {
					res += x
				}
				if argv.Round {
					res = math.Round(res)
				}
				fmt.Printf("Result: %f\n", res)
			},
		})

	// Don't round off the result
	cmd.Run([]string{"0.42", "3.14", "0.3"})

	// Round off the result
	cmd.Run([]string{"-round", "0.42", "3.14", "0.3"})

	// Output:
	// Result: 3.860000
	// Result: 4.000000
}

func ExampleCommand_third() {
	type NumArgs struct {
		Round bool      `flag:"round"`
		Args  []float64 `arg:"args"`
	}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:  "sum",
			Short: "adds floating point arguments",
			Long:  "<long description of how addition works>",
			Init: func() interface{} {
				var argv NumArgs

				// Make the default argument true
				argv.Round = true

				return &argv
			},
			Action: func(args interface{}) {
				argv, _ := args.(*NumArgs)
				res := 0.0
				for _, x := range argv.Args {
					res += x
				}
				if argv.Round {
					res = math.Round(res)
				}
				fmt.Printf("Result: %f\n", res)
			},
		})

	// This will round off the result
	cmd.Run([]string{"0.42", "3.14", "0.3"})

	// Turn off rounding
	cmd.Run([]string{"-round=false", "0.42", "3.14", "0.3"})

	// Output:
	// Result: 4.000000
	// Result: 3.860000
}

func ExampleNewMenu() {
	apply := func(init float64, op func(x, y float64) float64, args []float64) float64 {
		res := init
		for _, x := range args {
			res = op(res, x)
		}

		return res
	}
	add := func(x, y float64) float64 { return x + y }
	prod := func(x, y float64) float64 { return x * y }

	type NumCommand struct {
		Args []float64 `arg:"args" descr:"arguments"`
	}

	init := func() interface{} {
		return new(NumCommand)
	}
	doSum := func(args interface{}) {
		fmt.Printf("Result: %f\n", apply(0, add, args.(*NumCommand).Args))
	}
	doProd := func(args interface{}) {
		fmt.Printf("Result: %f\n", apply(1, prod, args.(*NumCommand).Args))
	}

	sum := cli.NewCommand(
		cli.CommandSpec{
			Name:   "+",
			Short:  "adds floating point arguments",
			Long:   "<long description of how addition works>",
			Init:   init,
			Action: doSum,
		})
	mult := cli.NewCommand(
		cli.CommandSpec{
			Name:   "*",
			Short:  "multiplies floating point arguments",
			Long:   "<long description of how multiplication works>",
			Init:   init,
			Action: doProd,
		})

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
	say := func(x string) func(interface{}) {
		return func(argv interface{}) { fmt.Println(x) }
	}
	menu := cli.NewMenu("menu", "", "",
		cli.NewMenu("foo", "", "",
			cli.NewCommand(cli.CommandSpec{Name: "x", Action: say("menu/foo/x")}),
			cli.NewCommand(cli.CommandSpec{Name: "y", Action: say("menu/foo/y")})),
		cli.NewMenu("bar", "", "",
			cli.NewCommand(cli.CommandSpec{Name: "x", Action: say("menu/bar/x")}),
			cli.NewCommand(cli.CommandSpec{Name: "y", Action: say("menu/bar/y")})),
		cli.NewCommand(cli.CommandSpec{Name: "baz", Action: say("menu/baz")}))

	menu.Run([]string{"baz"})      // will say menu/baz
	menu.Run([]string{"foo", "x"}) // will say menu/foo/x
	menu.Run([]string{"bar", "y"}) // will say menu/bar/y

	// Output:
	// menu/baz
	// menu/foo/x
	// menu/bar/y
}

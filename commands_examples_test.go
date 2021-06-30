package cli_test

import (
	"fmt"
	"math"

	"github.com/mailund/cli"
)

func ExampleCommand() {
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:  "first",
			Short: "my first command",
			Long:  "The first command I ever made",
		})

	cmd.Run([]string{})

	// Output:
}

func ExampleCommand_second() {
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "hello",
			Short:  "prints hello world",
			Action: func(_ interface{}) { fmt.Println("hello, world!") },
		})

	cmd.Run([]string{})

	// Output: hello, world!
}

func ExampleCommand_third() {
	type AddArgs struct {
		X int `pos:"x" descr:"first addition argument"`
		Y int `pos:"y" descr:"first addition argument"`
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

func ExampleCommand_fourth() {
	type NumArgs struct {
		Round bool      `flag:"round" descr:"should result be rounded to nearest integer"`
		Args  []float64 `pos:"args" descr:"numerical arguments to command"`
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
	cmd.Run([]string{"--round", "0.42", "3.14", "0.3"})

	// Output:
	// Result: 3.860000
	// Result: 4.000000
}

func ExampleCommand_fifth() {
	type NumArgs struct {
		Round bool      `flag:"round"`
		Args  []float64 `pos:"args"`
	}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:  "sum",
			Short: "adds floating point arguments",
			Long:  "<long description of how addition works>",
			Init: func() interface{} {
				return &NumArgs{Round: true}
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
	cmd.Run([]string{"--round=false", "0.42", "3.14", "0.3"})

	// Output:
	// Result: 4.000000
	// Result: 3.860000
}

func ExampleCommand_sixth() {
	type CalcArgs struct {
		X int `pos:"x" descr:"first addition argument"`
		Y int `pos:"y" descr:"first addition argument"`
	}

	add := cli.NewCommand(
		cli.CommandSpec{
			Name:  "add",
			Short: "adds two floating point arguments",
			Long:  "<long description of how addition works>",
			Init:  func() interface{} { return new(CalcArgs) },
			Action: func(args interface{}) {
				fmt.Printf("Result: %d\n", args.(*CalcArgs).X+args.(*CalcArgs).Y)
			},
		})

	mult := cli.NewCommand(
		cli.CommandSpec{
			Name:  "mult",
			Short: "multiplies two floating point arguments",
			Long:  "<long description of how multiplication works>",
			Init:  func() interface{} { return new(CalcArgs) },
			Action: func(args interface{}) {
				fmt.Printf("Result: %d\n", args.(*CalcArgs).X*args.(*CalcArgs).Y)
			},
		})

	calc := cli.NewCommand(
		cli.CommandSpec{
			Name:        "calc",
			Short:       "does calculations",
			Long:        "<long explanation of arithmetic>",
			Subcommands: []*cli.Command{add, mult},
		})

	calc.Run([]string{"add", "2", "3"})
	calc.Run([]string{"mult", "2", "3"})

	// Output:
	// Result: 5
	// Result: 6
}

func ExampleCommand_seventh() {
	type CalcArgs struct {
		X int `pos:"x" descr:"first addition argument"`
		Y int `pos:"y" descr:"first addition argument"`
	}

	add := cli.NewCommand(
		cli.CommandSpec{
			Name:  "add",
			Short: "adds two floating point arguments",
			Long:  "<long description of how addition works>",
			Init:  func() interface{} { return new(CalcArgs) },
			Action: func(args interface{}) {
				fmt.Printf("Result: %d\n", args.(*CalcArgs).X+args.(*CalcArgs).Y)
			},
		})

	mult := cli.NewCommand(
		cli.CommandSpec{
			Name:  "mult",
			Short: "multiplies two floating point arguments",
			Long:  "<long description of how multiplication works>",
			Init:  func() interface{} { return new(CalcArgs) },
			Action: func(args interface{}) {
				fmt.Printf("Result: %d\n", args.(*CalcArgs).X*args.(*CalcArgs).Y)
			},
		})

	calc := cli.NewCommand(
		cli.CommandSpec{
			Name:  "calc",
			Short: "does calculations",
			Long:  "<long explanation of arithmetic>",
			Action: func(_ interface{}) {
				fmt.Println("Hello from calc")
			},
			Subcommands: []*cli.Command{add, mult},
		})

	calc.Run([]string{"add", "2", "3"})
	calc.Run([]string{"mult", "2", "3"})

	// Output:
	// Hello from calc
	// Result: 5
	// Hello from calc
	// Result: 6
}

func ExampleNewMenu() {
	type CalcArgs struct {
		X int `pos:"x" descr:"first addition argument"`
		Y int `pos:"y" descr:"first addition argument"`
	}

	add := cli.NewCommand(
		cli.CommandSpec{
			Name:  "add",
			Short: "adds two floating point arguments",
			Long:  "<long description of how addition works>",
			Init:  func() interface{} { return new(CalcArgs) },
			Action: func(args interface{}) {
				fmt.Printf("Result: %d\n", args.(*CalcArgs).X+args.(*CalcArgs).Y)
			},
		})

	mult := cli.NewCommand(
		cli.CommandSpec{
			Name:  "mult",
			Short: "multiplies two floating point arguments",
			Long:  "<long description of how multiplication works>",
			Init:  func() interface{} { return new(CalcArgs) },
			Action: func(args interface{}) {
				fmt.Printf("Result: %d\n", args.(*CalcArgs).X*args.(*CalcArgs).Y)
			},
		})

	calc := cli.NewMenu(
		"calc", "does calculations", "<long explanation of arithmetic>",
		add, mult)

	calc.Run([]string{"add", "2", "3"})
	calc.Run([]string{"mult", "2", "3"})

	// Output:
	// Result: 5
	// Result: 6
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

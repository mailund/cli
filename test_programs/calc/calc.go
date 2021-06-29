package main

import (
	"fmt"
	"os"

	"github.com/mailund/cli"
)

type CalcArgs struct {
	X int `pos:"x" descr:"first addition argument"`
	Y int `pos:"y" descr:"first addition argument"`
}

func main() {
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

	calc.Run(os.Args[1:])
}

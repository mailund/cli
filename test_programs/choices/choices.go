package main

import (
	"fmt"
	"os"

	"github.com/mailund/cli"
)

type Args struct {
	A cli.Choice `flag:"" short:"a" descr:"optional choice"`
	B cli.Choice `pos:"b" descr:"mandatory choice"`
}

func Init() interface{} {
	return &Args{
		A: cli.Choice{Choice: "A", Options: []string{"A", "B", "C"}},
		B: cli.Choice{Choice: "B", Options: []string{"A", "B", "C"}},
	}
}

func Action(i interface{}) {
	args, _ := i.(*Args)
	fmt.Printf("Choice A was %s and choice B was %s\n", args.A.Choice, args.B.Choice)
}

func main() {
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "choices",
			Long:   "Demonstration of the difficult task of making choices.",
			Init:   Init,
			Action: Action,
		})

	cmd.Run(os.Args[1:])
}

package main

import (
	"fmt"
	"os"

	"github.com/mailund/cli"
)

type args struct {
	A cli.Choice `flag:"" short:"a" descr:"optional choice"`
	B cli.Choice `pos:"b" descr:"mandatory choice"`
}

func initChoices() interface{} {
	return &args{
		A: cli.Choice{Choice: "A", Options: []string{"A", "B", "C"}},
		B: cli.Choice{Choice: "B", Options: []string{"A", "B", "C"}},
	}
}

func choicesAction(i interface{}) {
	args, _ := i.(*args)
	fmt.Printf("Choice A was %s and choice B was %s\n", args.A.Choice, args.B.Choice)
}

func main() {
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "choices",
			Long:   "Demonstration of the difficult task of making choices.",
			Init:   initChoices,
			Action: choicesAction,
		})

	cmd.Run(os.Args[1:])
}

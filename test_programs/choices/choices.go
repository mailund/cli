package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailund/cli"
	"github.com/mailund/cli/interfaces"
)

type Choice struct {
	Choice  string
	Options []string
}

func (c *Choice) Set(x string) error {
	for _, v := range c.Options {
		if v == x {
			c.Choice = v
			return nil
		}
	}

	return interfaces.ParseErrorf(
		"%s is not a valid choice, must be in %s", x, "{"+strings.Join(c.Options, ",")+"}")
}

func (c *Choice) String() string {
	return c.Choice
}

func (c *Choice) ArgumentDescription(flag bool, descr string) string {
	if flag {
		// Don't modify the flag description (we handle it on the value)
		return descr
	}

	return descr + " (choose from " + "{" + strings.Join(c.Options, ",") + "})"
}

func (c *Choice) FlagValueDescription() string {
	return "{" + strings.Join(c.Options, ",") + "}"
}

type Args struct {
	A Choice `flag:"" short:"a" descr:"optional choice"`
	B Choice `pos:"b" descr:"mandatory choice"`
}

func Init() interface{} {
	return &Args{
		A: Choice{"A", []string{"A", "B", "C"}},
		B: Choice{"B", []string{"A", "B", "C"}},
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

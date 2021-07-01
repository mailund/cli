package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailund/cli"
	"github.com/mailund/cli/interfaces"
)

type Validator struct {
	validator func(string) bool
	x         string
}

func upperValidator(x string) bool { return strings.ToUpper(x) == x }
func lowerValidator(x string) bool { return strings.ToLower(x) == x }

func (val *Validator) setUpperValidator() {
	val.validator = upperValidator
}

func (val *Validator) setLowerValidator() {
	val.validator = lowerValidator
}

// Implementing Set(string) error means we can use a Validator as a
// positional argument
func (val *Validator) Set(x string) error {
	if val.validator != nil && !val.validator(x) {
		return interfaces.ParseErrorf("'%s' is not a valid string", x)
	}

	val.x = x

	return nil
}

type Args struct {
	Upper func()    `flag:"" short:"u" descr:"sets the validator to upper"`
	Lower func()    `flag:"" short:"l" descr:"sets the validator to lower"`
	Val   Validator `pos:"string" descr:"string we might do something to, if valid"`
}

// Now we have everything ready to set up the command.
func Init() interface{} {
	args := Args{}
	args.Upper = args.Val.setUpperValidator
	args.Lower = args.Val.setLowerValidator

	return &args
}

func Action(args interface{}) {
	fmt.Println("A valid string:", args.(*Args).Val.x)
}

func main() {
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "validator",
			Long:   "Will only accept valid strings",
			Init:   Init,
			Action: Action,
		})
	cmd.Run(os.Args[1:])
}

# cli

[![TravisCI Status](https://travis-ci.com/mailund/cli.svg?branch=main)](https://travis-ci.com/mailund/cli)
[![codecov](https://codecov.io/gh/mailund/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/mailund/cli)
[![Coverage Status](https://coveralls.io/repos/github/mailund/cli/badge.svg?branch=main)](https://coveralls.io/github/mailund/cli?branch=main)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4f8b3ef7896141b7ad4ace6f55d7ddc1)](https://www.codacy.com/manual/mailund/cli?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mailund/cli&amp;utm_campaign=Badge_Grade)

Command line parsing for go.

## Anatomy of command line tools as seen by cli

This package assumes that command line tools consists of a command, followed by flags, and then positional arguments:

```
> cmd -x -foo=42 a b c
```

The flags, `-x` and `-foo` are optional but with default values, while some of the positional arguments are required, except for potentially a variable number of arguments at the end of a command.

For many command line tools, there are also subcommands that looks like

```
> cmd -x -foo=42 subcmd -y a b c
```

where `cmd` is the command line tool, the first flags, `-x` and `-foo` are options to that, and then `subcmd` is a command in itself, getting its own flags, `-y` and arguments `a b c`.

Go's `flag` package handles flags, and does it very well, but it doesn't handle arguments or subcommands. There are other packages that implement support for subcommands, but not any that I much liked, so I implemented this module to get something more to my taste.

## Command line parameters

The package `params` in this module handles positional arguments. It is inspired by Go's `flag` module in how it handles them. To specify a set of arguments, you create an object of type `params.ParamSet`. It takes a name as its first argument, it is used when printing usage information if parsing failes, and a flag that determines whether parse errors should terminate the program or not. It is the same flag that `flag.NewFlags()` uses.

```go
p := params.NewParamSet("name", params.ExitOnError)
```

You probably always want to use `ExitOnError` for the flag. It will print usage information and terminate the program. The commands described below will always use this (although you can modify it if you want to). For command line tools, you generally want some helpful information printed, and then terminate, rather than trying to deal with incorrect parameters.

Once you have created the parameteter set, you install parameters of various type into it, the same way as you install flags with `flag`. For example, we can create a parameter set and install four variables, one for integers, one for boolean, one for floats and one for strings:

```go
// Variables we want to parse
var (
	i int
	b bool
)

p.IntVar(&i, "I", "an integer argument")
p.BoolVar(&b, "B", "a boolean argument")
f := p.Float("F", "a float argument")
s := p.String("S", "a string argument")
```

For each of the types, there is a `TypeVar()` function, where you provide a pointer to where you want the parsed value to go, or you can use a `Type()` function that will provide you with a pointer to where the value will go. The two string arguments are the name of the parameter and a short description, used when printing usage information. The order in which you install the parameters is the order in which you must provide them on the commandline.

When you have all your parameters, you can parse a commandline:

```go
p.Parse([]string{"42", "t", "3.14", "foo"})
// Now you have your arguments:
fmt.Printf("I=%d, B=%t, F=%f, S='%s'\n", i, b, *f, *s)
```

If parsing fails, and you were sensible enough to use `ExitOnError`, you won't continue after the parsing, so you can always assume that parameters were parsed correctly and are available in the variables they are associated with.

Somtimes you want more than one argument to a parameter. For a generic command line parser, it is impossible to figure out which arguments go where, unless it is the last parameter, so that is all we can handle here. You can install a variadic paramter, that will always be the last argument in the set (regardless of whether you install other parameters after, but please don't, since that can get confusing).

As an example of variadic paramters, consider a program where we want to calculate either the sum or product of a number of floating point values. Here, we can install a string parameter to pick the operation from, and then a variadic float parameter that will gobble up the rest of the arguments. You can provide a minimum number of arguments for a variadic parameter (but currently not a maximum), but in the example below I have set it to zero.

```go
p := params.NewParamSet("calc", params.ExitOnError)

// Variables we want to parse
var (
	op   string
	args []float64
)
p.StringVar(&op, "op", "operation to perform")
p.VariadicFloatVar(&args, "args", "arg [args...]", 0)

p.Parse([]string{"+", "0.42", "3.14", "0.3"})
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
```

The `op` parameter should either be '+' or '*', but the parser doesn't handle that yet. It is on my [TODO-list](https://github.com/mailund/cli/issues/1). The interesting part in the example, anyway, is the variadic parameter. It should be a slice of the right type, and when the arguments are parsed, the slice will have the parsed values.

Dispatching on an operator looks awfully lot like subcommands, so let us get to commands and subcommands...

## Commands

When you want a command line `Command`, you call `NewCommand` with four parameters. The first is the command's name, then you provide a short and a long description that is used when printing usage information, and finally a function that sets up the command. Needless to say, it is in the fourth parameter that the magic happens. Here is an example:

```go
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
```

The function should take a `cli.Command` as argument and provide a thunk (function without arguments) as its result. The object it is called with is the command we are creating, and it can use it to setup how the command should be executed (and change a few things in the command if necessary). In the example above, the callback function creates two arguments, `round` and `args`, adds `round` to the command's flags (those are of the type from Go's flag package) and adds `args` as a variadic positional argument. Then it returns a function that will add all the arguments in `args` when we call it, and round them off if the `round` flag is set.

The reason the design looks like this is that I want commands to call user-specified callback functions when we have parsed all arguments, potentially also after dispatching to sub-commands, and the callback should have access to the flags and parameters when they are parsed. There isn't any easy and type-safe method to wrap up parsed options and provide them to a generic function, but through this callback approach, the user can define flags and positional arguments in a closure, and thus the callback that the command gets will already know about them.

If you want to execute a command, you `Run()` it. Here's a call without the `round` flag:

```go
cmd.Run([]string{"0.42", "3.14", "0.3"})
```

and here is one where we include the flag

```go
cmd.Run([]string{"-round", "0.42", "3.14", "0.3"})
```

## Subcommands

Getting back to doing both addition and multiplication: here we have two different commands. One for each of the operations. We usually do not consider that different command line arguments, but why not? We could, and I will.

If you want to have sub-commands, you use the `NewMenu()` function. It creates a command wit subcommands. Here is a menu that lets you choose between executing a `+` or a `*` command:

```go
type operator = func(float64, float64) float64
type init_func = func(cmd *cli.Command) func()

apply := func(init float64, op operator, args []float64) float64 {
	res := init
	for _, x := range args {
		res = op(res, x)
	}
	return res
}
init_cmd := func(init float64, op operator) init_func {
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
	init_cmd(0, func(x, y float64) float64 { return x + y }))
mult := cli.NewCommand(
	"*", "multiplies floating point arguments",
	"<long description of how multiplication works>",
	init_cmd(1, func(x, y float64) float64 { return x * y }))

calc := cli.NewMenu("calc", "does calculations",
	"<long description of how math works>",
	sum, mult)
```

The first half of the example is just there because Go is a bit verbose when it comes to defining functions, and I want an apply function that generalises the two operators. I don't think Go has one, but if it does, please correct me. Anyway, I define an `apply` function that will do one operation over all floats in a slice, starting with some initial value. Then I write a function that creates a command initialisation callback by setting variadic arguments and then running `apply` with the `init` and `op` I give it. That has nothing to do with commands; I just didn't want a lot of replicated code like I had earlier.

Then, and this is where it is relevant for the `cli` module, I create a `sum` and `mult` command (named `+` and `*`), with the `init_cmd` function doing the setup. And with `NewMenu()` I put these in a command called `calc`.

A menu command takes a variable number of commands, the first of which must be the name of one of its subcommands (here `+` or `*`; it looks at the name we give the commands, not the variables). If you run `calc`, you specify which subcommand to via the first argument. This will call `sum`:

```go
calc.Run([]string{"+", "0.42", "3.14", "0.3"})
```

and this will call `mult`:

```go
calc.Run([]string{"*", "0.42", "3.14", "0.3"})
```

Any other first argument will be an error, and then `calc` will write some useful usage information and terminate the program.

You can nest menus arbitrarily deep. As the package parses a command line, a menu will always use its first positional parameter to dispatch to the next command, calling its `Run()` function with all arguments except the first (that was the commands name), so menus do not need to know where they are in the hierarhy. They will get a string slice of arguments, use the first to pick a subcommand, and then run the subcommand with the remaining. Once this procedure hits a non-menu command, it is run with the arguments that remain after moving into the menu hierarchy.

In this example, for example, different commands are picked based on the sequence of initial arguments:

```go
say := func(x string) func(*cli.Command) func() {
	// Sorry this is ugly, but I didn't want to write a lot of boiler plate
	// code below.
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
```

If the arguments get to an unknown command, the parsing will terminate with an error message, and the same happens if you end up at a menu command without a sub-command left to dispatch on.


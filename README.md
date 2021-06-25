# cli

[![TravisCI Status](https://travis-ci.com/mailund/cli.svg?branch=main)](https://travis-ci.com/mailund/cli)
[![codecov](https://codecov.io/gh/mailund/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/mailund/cli)
[![Coverage Status](https://coveralls.io/repos/github/mailund/cli/badge.svg?branch=main)](https://coveralls.io/github/mailund/cli?branch=main)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4f8b3ef7896141b7ad4ace6f55d7ddc1)](https://www.codacy.com/manual/mailund/cli?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mailund/cli&amp;utm_campaign=Badge_Grade)

Command line parsing for go.

## Anatomy of command line tools as seen by cli

This package assumes that command line tools consists of a command, followed by flags, and then positional arguments:

```sh
> cmd -x -foo=42 a b c
```

The flags, `-x` and `-foo` are optional but with default values, while some of the positional arguments are required, except for potentially a variable number of arguments at the end of a command.

For many command line tools, there are also subcommands that looks like

```sh
> cmd -x -foo=42 subcmd -y a b c
```

where `cmd` is the command line tool, the first flags, `-x` and `-foo` are options to that, and then `subcmd` is a command in itself, getting its own flags, `-y` and arguments `a b c`.

Go's `flag` package handles flags, and does it very well, but it doesn't handle arguments or subcommands. There are other packages that implement support for subcommands, but not any that I much liked, so I implemented this module to get something more to my taste.


## Commands

Commands are created with the `NewCommand` function from a command specification. The specification provides information about the name of the command, documentation as two strings, a short used when listing subcommands and a long used when showing usage of a specific command.

A command, that does absolutely nothing, but that you can run, could look like this:

```go
cmd := cli.NewCommand(
	cli.CommandSpec{
		Name:  "first",
		Short: "my first command",
		Long:  "The first command I ever made",
	})
```

If you want to run it, you call its `Run()` method with a slice of strings:

```go
cmd.Run([]string{})
```

This command doesn’t take any arguments, flags or positional, so you have to give it an empty slice.

If you want a command to do something when you call it, you must provide an `Action` in the specification. That is a function that takes an `interface{}` as its sole argument and doesn’t return anything. A “hello, world” command would look like this:

```go
cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "hello",
			Short:  "prints hello world",
			Action: func(_ interface{}) {
				fmt.Println("hello, world!")
			},
		})
```

If you run this command, `cmd.Run([]string{})`, it will print “hello, world!”.

The `interface{}` argument is where the action gets its arguments. To supply command arguments, you must define a type and one more function.

Let’s say we want a command that adds two numbers. For that, we need to arguments, and we will make them positional. We create a type for the arguments that looks like this:

```go
type AddArgs struct {
	X int `pos:"x" descr:"first addition argument"`
	Y int `pos:"y" descr:"first addition argument"`
}
```

The fields `X` and `Y` must be capitalised, because `cli` uses reflection to analyse the structure. The tags after the types is where we tell `cli` about the arguments. We make a field a positional argument using the tag `pos:"name"`, and if we want to give the parameter a description, we do so with `descr:"description"`. If you want a flag instead of a positional argument, you use `flag:"name"` instead of `pos:"name"`. You don’t need tags for all the fields in the structure you want to give your command, but those you want `cli` to parse, you do.

To set up arguments, you proved another function, `Init`. It should return an `interface{}`, and it is that `interface{}` that `cli` parses for tags and that it provides to `Action` when the command runs.

We could set up our command for adding two numbers like this:


```go
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
```

The function we provide to `Init` returns a pointer to a new `AddArgs`, and when `Action` is called, it can cast the `interface{}` argument to `*AddArgs` and get the positional values. If you run it with `add.Run([]string{"2", "2"})` it will print `Result: 4` as expected.

Here’s a variation that adds a flag as well:

```go
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
```

The type `NumArgs` has both a `flag:` and `pos:` tag, creating a boolean flag determining whether we should round the result of the calculations and a variadic (multiple argument) positional variable of type `[]float64`, which means that the command will be able to take any number of float arguments.

The `Init` function again just creates a `new(NumArgs)`, and the `Action` is a little more complicated, but only because it calculates the sum of multiple values. The way it handles the arguments is the same: it casts the argument to `*NumArgs` and from there it can get both the flag and the float positional arguments.

Run it with `cmd.Run([]string{"0.42", "3.14", "0.3"})` it and will print the sum (3.860000) and with `cmd.Run([]string{"-round", "0.42", "3.14", "0.3"})` and it will round the result (4.000000 — it is a silly example, so I still print it as a float…).

Flags have default values, so how do we deal with that? The short answer is that the default values for flags are the values the struct’s fields have when you give it to `cli`. That means that your `Init` function can set the default values simply by setting the struct’s fields before it returns it.

This is how we could make the `-round` flag default to `true`:

```go
cmd := cli.NewCommand(
	cli.CommandSpec{
		// The other arguments are the same as before…
		Init: func() interface{} {
			argv := new(NumArgs)
			argv.Round = true
			return argv
		},
		Action: func(args interface{}) {
			// the Action function is the same as before
		},
	})
```


Simply setting `argv.Round = true` before we return from `Init` will make `true` the default value for the flag.

You can nest commands to make subcommands. Say we want a tool that can do both addition and multiplication. We can create a command with two subcommands to achieve this. The straightforward way to do this is to give a command a `Subcommands` in its specification. It can look like this:

```go
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
```

Now `calc` has two subcommands, and you can invoke addition with `calc.Run([]string{"add", "2", "3"})` and multiplication with `calc.Run([]string{"must", "2", "3"})`. On the command line it would, of course, look like:

```sh
> calc add 2 3
Result: 5
> calc mult 2 3
Result: 6
```

The parent command, `calc` doesn’t have any action itself, and commands with subcommands do not have to. When a command has subcommands, `cli` automatically dispatches to the subcommands. If you provide an `Action`, however, it is called before the dispatch, and the dispatching is done after it complets. If we changed the `calc` command to this:

```go
calc := cli.NewCommand(
	cli.CommandSpec{
		Name:        "calc",
		Short:       "does calculations",
		Long:        "<long explanation of arithmetic>",
		Action:      func(_ interface{}) {
			fmt.Println("Hello from calc")
		},
		Subcommands: []*cli.Command{add, mult},
	})
```

the tool would print “Hello from calc” before dispatching:

```sh
> calc add 2 3
Hello from calc
Result: 5
> calc mult 2 3
Hello from calc
Result: 6
```

The syntax for subcommands, especially the `[]*cli.Command{…}` bit is somewhat cumbersome, and we usually do not associate arguments and actions to most commands that merely function as menus, so there is a convenience function that does the same thing:

```go
calc := cli.NewMenu(
	"calc", "does calculations",
	"<long explanation of arithmetic>",
	add, mult)
```

It takes the name, short and long parameters and then a variadic number of commands as arguments.

You can nest commands and subcommands arbitrarily deep:

```go
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
```

When you use subcommands, as you can see above, you don’t provide the name of the root command in the arguments. The root will usually be the executable, so if the executable for the example above is called `menu`, we would have:

```sh
> menu baz
menu/baz
> menu foo x
menu/foo/x
> menu bar y
menu/bar/y
```

so it all works out.

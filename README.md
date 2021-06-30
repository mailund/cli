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

If you want to make it a complete program, you just put the code in `main` and use `os.Args[1:]` as the arguments:

```go
package main

import (
  "fmt"
  "os"

  "github.com/mailund/cli"
)

func main() {
  cmd := cli.NewCommand(
    cli.CommandSpec{
      Name:   "hello",
      Long:   "Prints hello world",
      Action: func(_ interface{}) { fmt.Println("hello, world!") },
  })

  cmd.Run(os.Args[1:])
}
```

I changed the `Short` paramemter to `Long` here, since that is what `cli` uses to give full information about a command, rather than listing what subcommands do, but it wouldn't matter in this case. If `Long` isn't provided, it will use `Short`.

If you run the program without any arguments, you get `"hello, world!"` back. Assuming you compiled the program into the executable `hello`, you get:

```sh
> hello
hello, world!
```

If you call the program with the `-h` or `-help` flag, you get the following usage information:

```sh
> hello -h
Usage: hello [options]

Prints hello world

Options:
  -help
    show help for hello
```

If you call it with an unkown option, so anything but the help flag, or with any arguments, you get an error message and the same usage information. As we have specified the command, it doesn't take any arguments, so `cli` consider it an error if it gets any.

## Command arguments

The `interface{}` argument to the `Action` is where the action gets its arguments. To supply command arguments, you must define a type and one more function.

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
  Round bool      `flag:"round" short:"r" descr:"should result be rounded to nearest integer"`
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

Run it with `cmd.Run([]string{"0.42", "3.14", "0.3"})` it and will print the sum (3.860000) and with `cmd.Run([]string{"--round", "0.42", "3.14", "0.3"})` and it will round the result (4.000000 — it is a silly example, so I still print it as a float…).

The string after `flag:` in the tags defines the name of the flag. If it is a single letter, it can be used as both a long flag, `--f` and as a short flag `-f`, but if it is more than one letter, you can only use it as a long flag `--flag`. If you want both a long and a short flag, you can use the `short:` tag, as we did above. By adding `short:"r"`, we install two flags, the long `--round` and the short `-r`. Short flags can only have one letter, but you can combine them, so `-xyz` is equivalent to `-x -y -z`. You cannot combine long flags, but you can provide values to them with the syntax `--flag=value`, where short flags can only take values as `-f value`.

Flags have default values, so how do we deal with that? The short answer is that the default values for flags are the values the struct’s fields have when you give it to `cli`. That means that your `Init` function can set the default values simply by setting the struct’s fields before it returns it.

This is how we could make the `--round` flag default to `true`:

```go
cmd := cli.NewCommand(
  cli.CommandSpec{
    // The other arguments are the same as before…
    Init: func() interface{} {
      return &NumArgs{Round: true}
    },
    Action: func(args interface{}) {
      // the Action function is the same as before
    },
  })
```

Simply setting `argv.Round = true` before we return from `Init` will make `true` the default value for the flag.

You can use any of the types `string`, `bool`, `int`, `uint`, `int8`,  `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64` and `complex128` for flags and positional arguments. That's all the non-composite types except `uintptr` and unsafe pointers, which I wouldn't know how to handle generically...

Slice types with the same underlying types will be consider variadic arguments, and you can use those for positional arguments as long as there is only one variadic parameter per command, and provided that the command does not have sub-commands (see below).

Any type that implements the interface

```go
type PosValue interface {
  Set(string) error
}
```

for positional values or

```go
type FlagValue interface {
  String() string
  Set(string) error
}
```

as pointer-receivers will also work. Types that implement

```go
type VariadicValue interface {
	Set([]string) error
}
```

can be used as variadic parameters.


## Subcommands

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

Now `calc` has two subcommands, and you can invoke addition with `calc.Run([]string{"add", "2", "3"})` and multiplication with `calc.Run([]string{"must", "2", "3"})`.

Turn it into a complete program:

```go
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
```

compile it into an executable called `calc`, and then on the command line you can get information about how to use it using the `-h` flag

```sh
> calc -h
Usage: calc [options] cmd ...

<long explanation of arithmetic>

Options:
  -h,-help
    show help for calc

Arguments:
  cmd
	sub-command to call
  ...
	argument for sub-commands

Commands:
  add
	adds two floating point arguments
  mult
	multiplies two floating point arguments
```

You get the long description (which doesn't say much here, admittedly), then the options and arguments, since the `calc` command has subcommands those are `cmd` and `...`, and a list of the subcommands with their short description.

You can get help about the subcommands by providing those with the `-h` flag (or calling them incorrectly). To get help about `add`, we can use:

```sh
> calc add -h
Usage: add [options] x y

<long description of how addition works>

Options:
  -help
    show help for add

Arguments:
  x
	first addition argument
  y
	first addition argument
```

To invoke the commands, you do what experience has taught you and type the subcommand names after the main command:

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

## Callbacks

Instead of using parameters as values to set, you can install callback functions that will be called when the user provide a flag, or called on positional arguments. There are six types of functions you can use as callbacks, in different situations.

```go
type (
  NVCB  = func()
  NVCBI = func(interface{})

  NVCBE  = func() error
  NVCBIE = func(interface{}) error

  CB  = func(string) error
  CBI = func(string, interface{}) error

  VCB  = func([]string) error
  VCBI = func([]string, interface{}) error
)
```

The first four, `NVCB`, `NVCBI`, `NVCBE` and `NVCBIE` do not take any command-line arguments. They can only be used as flags. When they are, a flag, `-f` without arguments, will call the function. A `NVCB` and `NVCBE` flag will be called without arguments, of course, and a `NVCBI` or `NVCBIE` function will be called with the command line's argument structure, as returned from the `Init` function. The `-E` functions return an error and the other two do not. Callbacks that do not receive any input might not want to implement error handling, so `cli` support both boolean callbacks with and without an error return value. The functions that receive input should all return an `error` (although I might change that in the future).

The `CB` and `CBI` functions work with both flags and positional arguments. For flags, `--f=arg` or `-f arg`, the `arg` string is passed as the first argument, and for positional arguments, whatever argument is on the command line will be the callback's argument. They differ only in the second argument to functions of type `CBI`, which will be the command's argument-struct when the function is called (so values set before the function is called can be found in the struct, but flags and positional arguments that come after will not have been set yet).

The `VCB` and `VCBI` functions work the same as `CB` and `CBI` but take multiple arguments in the form of a string slice and can only be used for variadic parameters.

The example below shows callbacks and `Value` interfaces in action. We define a type, `Validator` that holds a string that must be validated according to some rules that are quite simple in this example, but you might be able to imagine something more complex.

We have two rules for validating the string in `Validator`, we might require that it is all uppercase or all lowercase, and we have functions to check this, and two methods for setting the rule. If we don't set a rule, then any string is valid.

If we implement the `PosValue` interface mentioned above, i.e. we implement a `Set(string) error` method, then we can use a `Validator` a positional argument, so we do that.

After that, we can make the type for command line arguments. We have a positional argument for the `Validator`, so it will get a positional argument, and we provide two flags to choose between the validation rules. If we don't provide a flag, we get the default (no checking), if we use `-u` we validate that the positional argument is all uppercase, and if we provide `-l` we validate that the argument is all lowercase.

The two callbacks have signature `func()` (the `BCB` type above), so they are valid types for `cli`. We put the methods from the validator the arguments holds in the fields. When we write `args.Val.setUpperValidator` we have already bound the receiver, so the method becomes a function that doesn't take any arguments, and thus we can use it as the field in the structure.

```go
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
  Upper func()    `flag:"u" descr:"sets the validator to upper"`
  Lower func()    `flag:"l" descr:"sets the validator to lower"`
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
```

Run it with `-h` to see a usage message:

```sh
> validate -h
Usage: validator [options] string

Will only accept valid strings

Options:
  -help show help for validator
  -l	  sets the validator to lower
  -u	  sets the validator to upper

Arguments:
  string
	string we might do something to, if valid
```

Without any flags, the program will accept any string:

```sh
> validate Foo
A valid string: Foo
```

but set a validator, and it will complain

```sh
validate -l Foo
Error parsing parameter string='Foo', 'Foo' is not a valid string.

> validate -l foo
A valid string: foo
```

## Protocols

Most of the behaviour of `cli` is handled through protocols and interfaces, making it relatively easy to extend.

Any type the implements `PosValue` as pointer-receiver can be used as a positional argument.

```go
type PosValue interface {
  Set(string) error // Should set the value from a string
}
```

The base types are wrapped in types that do this, and integers, for example, are handled like this:

```go
type IntValue int

func (val *IntValue) Set(x string) error {
  v, err := strconv.ParseInt(x, 0, strconv.IntSize)
  if err != nil {
    err = interfaces.ParseErrorf("argument \"%s\" cannot be parsed as int", x)
  } else {
    *val = IntValue(v)
  }
  return err
}
```

Flags also need a way to obtain the default value when printing usage help, so flags have the `FlagValue` interface:


```go
type FlagValue interface {
  String() string   // Should return a string representation of the value
  Set(string) error // Should set the value from a string
}
```

This is the same `Set(string) error` method as for positional arguments and a function for printing a value as a string. For integers, `cli`'s implementation is:

```go
func (val *IntValue) String() string {
	return strconv.Itoa(int(*val))
}
```

Implement your own type with the right methods, and you can use it as flags and positional arguments in `cli`. For example, if you want a type where you can select one of a small number of options, you could implement:

```go
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
    "%s is not a valid choice, must be in %s", x, "{" + strings.Join(c.Options, ",") + "}")
}

func (c *Choice) String() string {
  return c.Choice
}
```

This `Choice` type implements the `FlagValue`, so you can use it as both flags and positional arguments.


Variadic arguments implement the `VariadicValue` interface:

```go
type VariadicValue interface {
  Set([]string) error // Should set the value from a slice of strings
}
```

For integers, again, `cli` implements it as:

```go
type VariadicIntValue []int

func (vals *VariadicIntValue) Set(xs []string) error {
  *vals = make([]int, len(xs))

  for i, x := range xs {
    val, err := strconv.ParseInt(x, 0, strconv.IntSize)
    if err != nil {
      return interfaces.ParseErrorf("cannot parse '%s' as int", x)
    }

    (*vals)[i] = int(val)
	}

  return nil
}
```

For flags, there can be additional constraints. Some flags, for example, should not take arguments. If you don't want them to, then implement the `NoValueFlag` interface:

```go
type NoValueFlag interface {
  NoValueFlag() bool
}
```

This is how callbacks that do not take arguments inform the flag parser about that. This is how callbacks without arguments (but that may raise errors) are implemented:

```go
type FuncNoValue func() error         // Wrapping functions that do not take parameter arguments

func (f FuncNoValue) Set(_ string) error { return f() }
func (f FuncNoValue) String() string     { return "" }
func (f FuncNoValue) NoValueFlag() bool  { return true }
```

The first two methods ensure that we can use them as both flags and positional arguments, and the last tells the flags parser that it is an error to provide an argument, and a command line that looks like `--foo bar baz`, where `--foo` is such a flag, should not consider `bar` a value to be given to `--foo`.

Boolean flags are a little different. There, we usually want `--foo` to mean that the boolean value should be set to `true`, but we also want `--foo=true` and `--foo=false` to be valid. So, those flags have defaults. Well, all flags have defaults, but there is the default value when the flag is not used, and then the default value when we use it. For a boolean flag, usually the default if `false`, but `--foo` is interpreted as `--foo=true`, so `true` is the default string we should pass to the flag's `Set(string) error` method.

To handle those flags, we have the `DefaultValueFlag` interface:


```go
type DefaultValueFlag interface {
  Default() string
}
```

The `Default() string` method should return the default string to give to `Set(string) error` when the flag is used without an argument. For booleans, `cli` implement it as:

```go
func (val *BoolValue) Default() string { return "true" }
```

Flags that implement `DefaultValueFlag` can only take arguments as `--flag=arg`, since there is no way of knowing if `--flag foo` should interpret `foo` as a value for `--flag`, or if `--flag` should use its default and we should treat `foo` as a positional argument.

There are values that we want to validate as soon as we have linked struct fields to flags and parameters. Some data will crash the program when we try to parse a command line, and although we cannot capture this at compile time, when we use reflection to connect a struct with `cli`, we want to capture it early. We will eventually discover it when the parser reaches a point it cannot handle, of course, but it is better to check the data as soon as commands are connected, because then we catch it every time we run the program, rather than when commands are parsed.

The `Validator` interface gives hook for that:

```go
type Validator interface {
  Validate() error // Should return nil if everything is fine, or an error otherwise
}
```

Flags and parameters that implement the interface will have their `Validate() error` method called as soon as a command is created. Callbacks cannot be `nil` (that will certainly crash the program when we call them), and functions in `cli` check this using a `Validator` function:

```go
func (f FuncNoValue) Validate() error {
  if f == nil {
    return interfaces.SpecErrorf("callbacks cannot be nil")
  }

  return nil
}
```

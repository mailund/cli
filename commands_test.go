package cli_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/mailund/cli"
	"github.com/mailund/cli/internal/failure"
)

func TestNewCommand(t *testing.T) {
	_ = cli.NewCommand(cli.CommandSpec{
		Name:  "name",
		Short: "short",
		Long:  "long",
	})
}

func TestNewCommandError(t *testing.T) {
	type Invalid struct {
		X uintptr `flag:"x"`
	}

	_, err := cli.NewCommandError(cli.CommandSpec{
		Init: func() interface{} { return new(Invalid) },
	})

	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestNewCommandUsage(t *testing.T) { //nolint:funlen // Only too long because of usage strings
	type argv struct {
		Foo string `flag:"foo" descr:"set foo"`
		Bar string `pos:"bar" descr:"a bar"`
	}

	cmd := cli.NewCommand(cli.CommandSpec{
		Name:  "name",
		Short: "short",
		Long:  "long",
		Init:  func() interface{} { return new(argv) }})

	builder := new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.Usage()

	msg := builder.String()
	expected := `Usage: name [flags] bar
        
        
long
	
Flags:
  -h,--help
	show help for name
  --foo value
	set foo
	
Arguments:
  bar
	a bar

`
	space := regexp.MustCompile(`\s+`)
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")

	if msg != expected {
		t.Errorf("Expected usage message %s but got %s\n", expected, msg)
	}

	cmd = cli.NewCommand(cli.CommandSpec{
		Name:  "name",
		Short: "short",
		Init:  func() interface{} { return new(argv) }})

	builder = new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.Usage()

	msg = builder.String()
	expected = `Usage: name [flags] bar
        
        
short
	
Flags:
  -h,--help
	show help for name
  --foo value
	set foo
	
Arguments:
  bar
	a bar

`
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")

	if msg != expected {
		t.Errorf("Expected usage message %s but got %s\n", expected, msg)
	}

	failure.Success = func() {}

	builder = new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.Run([]string{"--help", "bar"})

	msg = space.ReplaceAllString(builder.String(), " ")
	if msg != expected {
		t.Errorf("Expected usage message %s but got %s\n", expected, msg)
	}
}

func TestNewCommandUsage2(t *testing.T) {
	type Argv struct {
		Foo string `flag:"foo" short:"f" descr:"foo"`
		Bar bool   `flag:"bar" descr:"bar"`
		Baz func() `flag:"baz" descr:"baz"`
	}

	argv := &Argv{
		Baz: func() {},
	}

	cmd := cli.NewCommand(cli.CommandSpec{
		Name:  "name",
		Short: "short",
		Long:  "long",
		Init:  func() interface{} { return argv },
	})

	builder := new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.Usage()

	expected := `Usage: name [flags]

	long


	Flags:
	  -h,--help
		show help for name
	  -f,--foo value
		foo
	  --bar [value] (no value = true)
		bar (default false)
	  --baz
		baz
`

	msg := builder.String()

	space := regexp.MustCompile(`\s+`)
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")

	if msg != expected {
		t.Errorf("unexpected: %s", msg)
		t.Errorf("unexpected: %s", expected)
	}
}

func TestCommandCalled(t *testing.T) {
	called := false
	cmd := cli.NewCommand(cli.CommandSpec{
		Action: func(argv interface{}) { called = true },
	})
	cmd.Run([]string{})

	if !called {
		t.Fatal("Command x wasn't called")
	}
}

func TestOption(t *testing.T) {
	called := false

	type argv struct {
		X int `flag:"xx" short:"x"`
	}

	var parsed argv

	init := func() interface{} { return &parsed }
	cmd := cli.NewCommand(cli.CommandSpec{
		Name:   "foo",
		Short:  "does foo",
		Long:   "a really good foo",
		Init:   init,
		Action: func(_ interface{}) { called = true }})
	cmd.SetErrorFlag(failure.ContinueOnError)

	// We won't call Failure() here, because the error is handled
	// by flag.FlagSet(), but we get the error from it.
	builder := new(strings.Builder)
	cmd.SetOutput(builder)

	cmd.Run([]string{"-x", "foo"}) // wrong type for int

	if called {
		t.Error("The command shouldn't be called")
	}

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, `error parsing flag -x: argument "foo" cannot be parsed as int.`) {
		t.Errorf("Unexpected error msg: %s", errmsg)
	}

	cmd.Run([]string{"-x", "42"}) // correct type for int

	if !called {
		t.Error("The command should be called now")
	}

	if parsed.X != 42 {
		t.Error("The option wasn't set correctly")
	}

	// the long flag should work as well...
	called = false
	parsed.X = 0

	cmd.Run([]string{"--xx", "42"}) // correct type for int

	if !called {
		t.Error("The command should be called now")
	}

	if parsed.X != 42 {
		t.Error("The option wasn't set correctly")
	}
}

func TestParam(t *testing.T) {
	type Testargs struct {
		X int `pos:"x" descr:"an integer"`
	}

	argX := 0
	failed := false
	called := false
	init := func() interface{} {
		var a Testargs
		a.X = 43

		return &a
	}
	action := func(a interface{}) {
		called = true
		argX = a.(*Testargs).X
	}

	failure.Failure = func() { failed = true }

	cmd := cli.NewCommand(cli.CommandSpec{
		Name:   "foo",
		Short:  "does foo",
		Long:   "a really good foo",
		Init:   init,
		Action: action,
	})

	builder := new(strings.Builder)
	cmd.SetOutput(builder)

	cmd.Run([]string{"foo"}) // wrong type for int

	if !failed {
		t.Error("The command should have failed")
	}

	if called {
		t.Error("The command shouldn't be called")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, `Error parsing parameter x='foo'`) {
		t.Errorf("Unexpected error msg: %s", errmsg)
	}

	failed = false

	cmd.Run([]string{"42"}) // correct type for int

	if failed {
		t.Error("Didn't expect a failure this time")
	}

	if !called {
		t.Error("The command should be called now")
	}

	if argX != 42 {
		t.Error("The option wasn't set correctly")
	}
}

func makeMenu() (xCalled, yCalled *bool, cmd *cli.Command) {
	xCalled, yCalled = new(bool), new(bool)

	actionFunc := func(b *bool) func(argv interface{}) {
		return func(argv interface{}) {
			*b = true
		}
	}
	cmd = cli.NewMenu("menu", "", "Dispatch to subcommands",
		cli.NewCommand(cli.CommandSpec{Name: "x", Short: "do x", Action: actionFunc(xCalled)}),
		cli.NewCommand(cli.CommandSpec{Name: "y", Short: "do y", Action: actionFunc(yCalled)}),
	)

	return xCalled, yCalled, cmd
}

func TestMenu(t *testing.T) {
	xCalled, yCalled, menu := makeMenu()

	menu.Run([]string{"x"})

	if !*xCalled {
		t.Error("Command x wasn't called")
	}

	if *yCalled {
		t.Error("y was called too soon")
	}

	menu.Run([]string{"y"})

	if !*yCalled {
		t.Error("Command y wasn't called")
	}
}

func TestMenuUsage(t *testing.T) {
	_, _, menu := makeMenu()

	builder := new(strings.Builder)
	menu.SetOutput(builder)
	cli.DefaultUsage(menu)

	defaultUsage := builder.String()
	builder = new(strings.Builder)

	menu.SetOutput(builder)
	menu.Usage()

	var menuUsage = builder.String()

	if !strings.HasPrefix(menuUsage, defaultUsage) {
		t.Error("Expected menu usage to start with default usage")
	}

	if !strings.HasSuffix(menuUsage, "Commands:\n  x\n\tdo x\n  y\n\tdo y\n\n") {
		t.Error("Expected commands at the end of menu usage")
		t.Errorf("The usage output was: %s", menuUsage)
	}
}

func TestMenuFailure(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }
	_, _, menu := makeMenu()

	builder := new(strings.Builder)
	menu.SetOutput(builder)
	menu.SetErrorFlag(failure.ContinueOnError)
	menu.Run([]string{"z"})

	if !failed {
		t.Error("Expected command to fail")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, "'z' is not a valid command for menu.") {
		t.Errorf("(1) Expected different error message than %s\n", errmsg)
	}

	builder = new(strings.Builder)
	menu.SetOutput(builder)
	menu.Run([]string{"x", "--foo"})

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "flag provided but not defined: --foo.") {
		t.Errorf("(2) Expected different error message than %s\n", errmsg)
	}
}

func TestCommandPanic(t *testing.T) {
	defer func() { _ = recover() }()

	type Invalid struct {
		X uintptr `pos:"invalid type"`
	}

	_ = cli.NewCommand(cli.CommandSpec{
		Init: func() interface{} { return new(Invalid) },
	})

	t.Error("The code did not panic")
}

func TestPositionalArgsWithSubcommands(t *testing.T) {
	type Args struct {
		X string `pos:"x" descr:"an x that does x"`
	}

	subcmdCalled := false
	cmdCalled := false
	x := ""
	subcmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "subcmd",
			Action: func(_ interface{}) { subcmdCalled = true }})
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:        "cmd",
			Subcommands: []*cli.Command{subcmd},

			Init: func() interface{} { return new(Args) },
			Action: func(a interface{}) {
				cmdCalled = true
				x = a.(*Args).X
			}})

	cmd.Run([]string{"foo", "subcmd"})

	if !subcmdCalled {
		t.Error("subcmd should have been called")
	}

	if !cmdCalled {
		t.Error("cmd should have been called")
	}

	if x != "foo" {
		t.Errorf("x was set to %s, we expected foo", x)
	}
}

func TestVariadicArgsWithSubcommands(t *testing.T) {
	type Args struct {
		X []string `pos:"x" descr:"an x that does x"`
	}

	subcmd := cli.NewCommand(cli.CommandSpec{Name: "subcmd"})
	_, err := cli.NewCommandError(
		cli.CommandSpec{
			Init:        func() interface{} { return new(Args) },
			Subcommands: []*cli.Command{subcmd}})

	if err == nil {
		t.Fatal("Expected an error")
	}

	if err.Error() != "a command with subcommands cannot have variadic parameters" {
		t.Errorf("Unexpected error message: '%s'\n", err.Error())
	}
}

type ValueInterface struct {
	x string
	i int
}

func (val *ValueInterface) Set(s string) error {
	val.x = s
	val.i = len(s)

	return nil
}

func (val *ValueInterface) String() string {
	return val.x
}

type VariadicValueInterface struct {
	x []string
	i int
}

func (val *VariadicValueInterface) Set(xs []string) error {
	val.x = xs
	val.i = len(xs)

	return nil
}

func TestValueInterface(t *testing.T) {
	var args = struct {
		Flag ValueInterface         `flag:"value"`
		Pos  ValueInterface         `pos:"value"`
		Var  VariadicValueInterface `pos:"var"`
	}{}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Init: func() interface{} { return &args },
		})

	cmd.Run([]string{"--value=foobar", "qux", "the", "rest", "of", "the", "args"})

	if args.Flag.x != "foobar" || args.Flag.i != 6 {
		t.Error("The flag wasn't set correctly")
	}

	if args.Pos.x != "qux" || args.Pos.i != 3 {
		t.Error("The positional wasn't set correctly")
	}

	if args.Var.i != 5 || !reflect.DeepEqual(args.Var.x, []string{"the", "rest", "of", "the", "args"}) {
		t.Error("Variadic argument wasn't set correctly")
	}
}

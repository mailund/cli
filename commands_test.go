package cli_test

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/mailund/cli"
	"github.com/mailund/cli/internal/failure"
)

func TestNewCommand(t *testing.T) {
	cli.NewCommand("name", "short", "long",
		func(cmd *cli.Command) func() { return func() {} })
}

func TestNewCommandUsage(t *testing.T) {
	init := func(cmd *cli.Command) func() {
		cmd.Flags.Bool("foo", false, "set foo")
		cmd.Params.String("bar", "a bar")
		return func() {}
	}
	cmd := cli.NewCommand("name", "short", "long", init)

	builder := new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.Usage()
	msg := builder.String()
	expected := `Usage: name [options] bar
        
        
long
	
Options:
  -foo
	set foo
  -help
	show help for name
	
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

	cmd = cli.NewCommand("name", "short", "", init)

	builder = new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.Usage()
	msg = builder.String()
	expected = `Usage: name [options] bar
        
        
short
	
Options:
  -foo
	set foo
  -help
	show help for name
	
Arguments:
  bar
	a bar

`
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")
	if msg != expected {
		t.Errorf("Expected usage message %s but got %s\n", expected, msg)
	}
}

func TestCommandCalled(t *testing.T) {
	called := false
	init := func(cmd *cli.Command) func() { return func() { called = true } }
	cmd := cli.NewCommand("", "", "", init)
	cmd.Run([]string{})
	if !called {
		t.Fatal("Command x wasn't called")
	}
}

func TestOption(t *testing.T) {
	called := false
	argX := -1
	init := func(cmd *cli.Command) func() {
		x := cmd.Flags.Int("x", 42, "an integer")
		return func() { called = true; argX = *x }
	}
	cmd := cli.NewCommand("foo", "does foo", "a really good foo", init)
	cmd.Flags.Init("foo", flag.ContinueOnError) // so we don't terminate

	// We won't call Failure() here, because the error is handled
	// by flag.FlagSet(), but we get the error from it.
	builder := new(strings.Builder)
	cmd.SetOutput(builder)

	cmd.Run([]string{"-x", "foo"}) // wrong type for int
	if called {
		t.Error("The command shouldn't be called")
	}
	errmsg := builder.String()
	if !strings.HasPrefix(errmsg, `invalid value "foo" for flag -x: parse error`) {
		t.Errorf("Unexpected error msg: %s", errmsg)
	}

	cmd.Run([]string{"-x", "42"}) // correct type for int
	if !called {
		t.Error("The command should be called now")
	}
	if argX != 42 {
		t.Error("The option wasn't set correctly")
	}
}

func TestParam(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	called := false
	argX := -1
	init := func(cmd *cli.Command) func() {
		x := cmd.Params.Int("x", "an integer")
		return func() { called = true; argX = *x }
	}
	cmd := cli.NewCommand("foo", "does foo", "a really good foo", init)

	builder := new(strings.Builder)
	cmd.SetOutput(builder)

	cmd.Run([]string{"foo"}) // wrong type for int
	if !failed {
		t.Error("The command should have failed")
	}
	if called {
		t.Error("The command shouldn't be called")
	}
	errmsg := builder.String()
	if !strings.HasPrefix(errmsg, `Error parsing parameter x='foo'`) {
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

func makeMenu() (*bool, *bool, *cli.Command) {
	xCalled, yCalled := false, false
	initFunc := func(b *bool) func(*cli.Command) func() {
		return func(cmd *cli.Command) func() {
			return func() { *b = true }
		}
	}
	menu := cli.NewMenu("menu", "", "Dispatch to subcommands",
		cli.NewCommand("x", "do x", "", initFunc(&xCalled)),
		cli.NewCommand("y", "do y", "", initFunc(&yCalled)),
	)
	return &xCalled, &yCalled, menu
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
	menuUsage := builder.String()

	if !strings.HasPrefix(menuUsage, defaultUsage) {
		t.Error("Expected menu usage to start with default usage")
	}
	if !strings.HasSuffix(menuUsage, "Commands:\n  x\n\tdo x\n  y\n\tdo y\n\n") {
		t.Error("Expected commands at the end of menu usage")
		fmt.Println("The usage output was:")
		fmt.Println(menuUsage)
	}
}

func TestMenuFailure(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }
	_, _, menu := makeMenu()

	builder := new(strings.Builder)
	menu.SetOutput(builder)
	menu.Run([]string{"z"})
	errmsg := builder.String()

	if !failed {
		t.Error("Expected command to fail")
	}
	if !strings.HasPrefix(errmsg, "'z' is not a valid command for menu.") {
		t.Errorf("Expected different error message than %s\n", errmsg)
	}
}

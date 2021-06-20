package cli_test

import (
	"flag"
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
	arg_x := -1
	init := func(cmd *cli.Command) func() {
		x := cmd.Flags.Int("x", 42, "an integer")
		return func() { called = true; arg_x = *x }
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
	if arg_x != 42 {
		t.Error("The option wasn't set correctly")
	}
}

func TestParam(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	called := false
	arg_x := -1
	init := func(cmd *cli.Command) func() {
		x := cmd.Params.Int("x", "an integer")
		return func() { called = true; arg_x = *x }
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
	if arg_x != 42 {
		t.Error("The option wasn't set correctly")
	}
}

func TestMenu(t *testing.T) {
	x_called, y_called := false, false
	init_func := func(b *bool) func(*cli.Command) func() {
		return func(cmd *cli.Command) func() {
			return func() { *b = true }
		}
	}
	menu := cli.NewMenu("", "", "",
		cli.NewCommand("x", "", "", init_func(&x_called)),
		cli.NewCommand("y", "", "", init_func(&y_called)),
	)

	menu.Run([]string{"x"})
	if !x_called {
		t.Error("Command x wasn't called")
	}
	if y_called {
		t.Error("y was called too soon")
	}
	menu.Run([]string{"y"})
	if !y_called {
		t.Error("Command y wasn't called")
	}

}

package params_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/internal/failure"
	"github.com/mailund/cli/params"
)

func TestMakeParamSet(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	if p.Name != "test" {
		t.Fatalf("New params set got the wrong name")
	}
}

func TestShortUsage(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	if p.ShortUsage() != "" {
		t.Errorf(`Short usage of empty paramset should be ""`)
	}
	p.String("foo", "")
	if p.ShortUsage() != "foo" {
		t.Errorf(`Short usage of with one "foo" parameter should be "foo"`)
	}
	p.String("bar", "")
	if p.ShortUsage() != "foo bar" {
		t.Errorf(`Short usage of with parameters "foo" and "bar" should be "foo bar"`)
	}
	p.VariadicString("...", "", 0)
	if p.ShortUsage() != "foo bar ..." {
		t.Errorf(`Short usage of with parameters "foo", "bar" and moreshould be "foo bar ..."`)
	}
}

func TestStringParam(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	var x string
	p.StringVar(&x, "x", "")
	var y *string = p.String("y", "")
	p.Parse([]string{"foo", "bar"})
	if x != "foo" {
		t.Errorf(`Expected var x to hold "foo"`)
	}
	if *y != "bar" {
		t.Errorf(`Expected var *y to hold "bar"`)
	}
}

func TestPrintDefault(t *testing.T) {
	builder := new(strings.Builder)
	p := params.NewParamSet("test", params.ExitOnError)
	p.SetOutput(builder)

	var x string
	p.StringVar(&x, "x", "an x")
	p.String("y", "a y")

	p.PrintDefaults()
	res := builder.String()
	expected := `Arguments:
  x
	an x
  y
	a y
`
	if res != expected {
		t.Fatalf("PrintDefault gave us \n----\n%s\n----\n\nbut we exected\n\n----\n%s\n----",
			res, expected)
	}
}

func TestPrintDefaultMore(t *testing.T) {
	builder := new(strings.Builder)
	p := params.NewParamSet("test", params.ExitOnError)
	p.SetOutput(builder)

	var x string
	p.StringVar(&x, "x", "an x")
	p.String("y", "a y")
	p.VariadicString("...", "arguments for ...", 0)

	p.PrintDefaults()
	res := builder.String()
	expected := `Arguments:
  x
	an x
  y
	a y
  ...
	arguments for ...
`
	if res != expected {
		t.Fatalf("PrintDefault gave us \n----\n%s\n----\n\nbut we exected\n\n----\n%s\n----",
			res, expected)
	}
}

func TestParseMore(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)

	p.String("x", "")
	rest := p.VariadicString("...", "arguments for ...", 0)

	args := []string{"for x", "y", "z"}
	p.Parse(args)
	if !reflect.DeepEqual(*rest, []string{"y", "z"}) {
		t.Fatalf("The parser ate more than it should!")
	}
}

func TestFailure(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)

	p.String("x", "")

	var failed = false
	failure.Failure = func() { failed = true }
	builder := new(strings.Builder)
	p.SetOutput(builder)
	err := p.Parse([]string{})
	errmsg := builder.String()
	if !failed {
		t.Errorf("We expected a parse failure")
	}
	if !strings.HasPrefix(errmsg, "Too few arguments") {
		t.Errorf("We got an unexpected error message: %s", errmsg)
	}
	if err == nil {
		t.Fatal("expected an error from the parse error")
	}
	if err.Error() != "too few arguments" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	failed = false
	builder = new(strings.Builder)
	p.SetOutput(builder)
	err = p.Parse([]string{"x", "y"})
	errmsg = builder.String()
	if !failed {
		t.Errorf("We expected a parse failure")
	}
	if !strings.HasPrefix(errmsg, "Too many arguments") {
		t.Errorf("We got an unexpected error message: %s", errmsg)
	}
	if err == nil {
		t.Fatal("expected an error from the parse error")
	}
	if err.Error() != "too many arguments" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	// Add a variadic that wants at least one argument
	p.VariadicString("y", "", 1)
	failed = false
	builder = new(strings.Builder)
	p.SetOutput(builder)
	err = p.Parse([]string{"x"}) // too few
	errmsg = builder.String()
	if !failed {
		t.Errorf("We expected a parse failure")
	}
	if !strings.HasPrefix(errmsg, "Too few arguments") {
		t.Errorf("We got an unexpected error message: %s", errmsg)
	}
	if err == nil {
		t.Fatal("expected an error from the parse error")
	}
	if err.Error() != "too few arguments" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	failed = false
	err = p.Parse([]string{"x", "y"}) // now we have the minimum
	if failed {
		t.Errorf("We shouldn't fail")
	}
	if err != nil {
		t.Fatal("We shouldn't get an error now")
	}

	err = p.Parse([]string{"x", "y", "z", "w"}) // and we can have more
	if failed {
		t.Errorf("We shouldn't fail")
	}
	if err != nil {
		t.Fatal("We shouldn't get an error now")
	}
}

func TestFailureContinue(t *testing.T) {
	p := params.NewParamSet("test", params.ContinueOnError)

	builder := new(strings.Builder)
	p.SetOutput(builder)

	p.String("x", "")
	var failed = false
	failure.Failure = func() { failed = true }
	err := p.Parse([]string{})
	if err == nil {
		t.Fatalf("expected an error from the parse error")
	}
	errmsg := builder.String()
	if failed {
		t.Fatal("We shouldn't terminate now")
	}
	if errmsg != "" {
		t.Fatalf("We shouldn't have an error message (but have %s)", errmsg)
	}
}

func TestFuncCallback(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	var called = false
	f := func(arg string) error {
		if arg != "arg" {
			t.Errorf("Got unexpected argument for callback")
		}
		called = true
		return nil
	}
	p.Func("foo", "", f)
	p.Parse([]string{"arg"})
	if !called {
		t.Fatalf("Callback was not called")
	}
}

func TestFuncCallbackError(t *testing.T) {
	var failed = false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", params.ExitOnError)
	builder := new(strings.Builder)
	p.SetOutput(builder)

	p.Func("foo", "",
		func(arg string) error {
			return errors.New("foo failed to bar")
		})
	err := p.Parse([]string{"arg"})
	if err == nil {
		t.Fatalf("Expected an error")
	}
	errmsg := builder.String()

	if !failed {
		t.Errorf("Expected parse error to trigger a failure")
	}
	if !strings.HasSuffix(errmsg, "foo failed to bar\n") {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}

func TestVariadicFuncError(t *testing.T) {
	var failed = false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", params.ExitOnError)
	builder := new(strings.Builder)
	p.SetOutput(builder)

	p.VariadicFunc("foo", "", 0,
		func(args []string) error {
			return errors.New("foo failed to bar")
		})
	err := p.Parse([]string{"arg"})
	if err == nil {
		t.Fatalf("Expected an error")
	}
	errmsg := builder.String()

	if !failed {
		t.Errorf("Expected parse error to trigger a failure")
	}
	if !strings.HasSuffix(errmsg, "foo failed to bar\n") {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}
func TestVariadicStrings(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	x := p.String("x", "")
	res := p.VariadicString("x [x...]", "", 0)
	args := []string{"x", "y", "z"}
	p.Parse(args)
	if *x != "x" {
		t.Errorf("Argument x should be x, is %s", *x)
	}
	if !reflect.DeepEqual(args[1:], *res) {
		t.Errorf("Variadic argument got %v, expected [y, z]", *res)
	}
}

func TestInt(t *testing.T) {
	var failed = false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", params.ExitOnError)
	ip := p.Int("i", "int")
	p.Parse([]string{"42"})
	if *ip != 42 {
		t.Errorf("Parse error, ip is %d", *ip)
	}

	builder := new(strings.Builder)
	p.SetOutput(builder)
	err := p.Parse([]string{"foo"})
	if err == nil {
		t.Error("Expected an error")
	}
	if !failed {
		t.Error("Expected to fail")
	}
	errmsg := builder.String()
	if !strings.HasPrefix(errmsg, `Error parsing parameter i='foo'`) {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}
